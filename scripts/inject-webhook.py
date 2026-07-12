#!/usr/bin/env python3
"""Inject conversion webhook and fix v1beta1 singleton-list schemas at build time.

Without spec.conversion.strategy: Webhook in the CRD YAML, Crossplane's
APIEstablisher skips CA bundle injection — the CRD stays at strategy: None
and the singleton-list conversion webhook never activates.

This script mirrors the behavior of `make kustomize-crds` which applies
package/kustomize/webhook.yaml via kustomize to inject the webhook config.

It also patches v1beta1 CRD schemas to use objects instead of arrays for
singleton-list fields. The AddSingletonListConversion generates v1beta1
with arrays and v1beta2 with objects, which breaks Kubernetes Server-Side
Apply (SSA). SSA's field manager can't reconcile arrays in v1beta1 with
objects in v1beta2 — it crashes with "expected list, got map".

Making v1beta1 use objects (matching v1beta2) eliminates the type mismatch
and SSA processes both versions identically.
"""

import sys
import pathlib
import yaml

WEBHOOK = {
    "strategy": "Webhook",
    "webhook": {
        "clientConfig": {
            "service": {
                "name": "",
                "namespace": "",
                "path": "/convert",
            }
        },
        "conversionReviewVersions": ["v1"],
    },
}

# Resources whose v1beta1 schemas need singleton-list array→object patching.
# These resources use AddSingletonListConversion which creates incompatible
# types between v1beta1 (arrays) and v1beta2 (objects), breaking SSA.
SINGLETON_LIST_RESOURCES = {
    "applications.azuread.upbound.io": {
        "plural": "applications",
        "fields": [
            "api",
            "featureTags",
            "optionalClaims",
            "password",
            "publicClient",
            "singlePageApplication",
            "timeouts",
            "web",
        ],
    },
}


def patch_for_provider(parent: dict, fields: list[str]) -> bool:
    """Patch array→object in forProvider and initProvider for listed fields."""
    changed = False
    for section in ("forProvider", "initProvider"):
        props = (
            parent.get("properties", {})
            .get(section, {})
            .get("properties", {})
        )
        for field in fields:
            if field not in props:
                continue
            f = props[field]
            if f.get("type") != "array":
                continue
            items = f.get("items")
            if not items or items.get("type") != "object":
                continue
            # Convert array-of-object → object
            f["type"] = "object"
            f["properties"] = items.get("properties", {})
            # Keep description if present
            if "description" not in f and "description" in items:
                f["description"] = items["description"]
            f.pop("items", None)
            f.pop("x-kubernetes-list-type", None)
            f.pop("x-kubernetes-list-map-keys", None)
            changed = True
    return changed


def patch_v1beta1_schema(doc: dict) -> bool:
    """Patch v1beta1 CRD schema to use objects for singleton-list fields."""
    group = doc.get("spec", {}).get("group", "")
    plural = doc.get("spec", {}).get("names", {}).get("plural", "")
    key = f"{group}"

    if key not in SINGLETON_LIST_RESOURCES:
        return False
    cfg = SINGLETON_LIST_RESOURCES[key]
    if plural != cfg["plural"]:
        return False

    changed = False
    for version in doc.get("spec", {}).get("versions", []):
        if version.get("name") != "v1beta1":
            continue
        schema = version.get("schema", {}).get("openAPIV3Schema", {})
        spec = schema.get("properties", {}).get("spec", {})

        if patch_for_provider(spec, cfg["fields"]):
            changed = True

    return changed


def main(crd_dir: str) -> None:
    crd_path = pathlib.Path(crd_dir)
    if not crd_path.is_dir():
        print(f"ERROR: {crd_dir} is not a directory", file=sys.stderr)
        sys.exit(1)

    webhook_count = 0
    schema_count = 0
    for f in sorted(crd_path.glob("*.yaml")):
        try:
            doc = yaml.safe_load(f.read_text())
        except Exception as e:
            print(f"  SKIP  {f.name}: cannot parse YAML ({e})")
            continue

        if doc is None or doc.get("kind") != "CustomResourceDefinition":
            continue

        changed = False

        # Inject conversion webhook
        spec = doc.setdefault("spec", {})
        spec["conversion"] = WEBHOOK
        webhook_count += 1
        changed = True

        # Patch v1beta1 singleton-list schemas
        if patch_v1beta1_schema(doc):
            schema_count += 1
            changed = True

        if changed:
            f.write_text(
                yaml.dump(doc, default_flow_style=False, sort_keys=False)
            )
        print(f"  {'v1beta1+webhook' if schema_count > 0 and webhook_count > 0 else 'webhook' if webhook_count > 0 else 'schema'} → {f.name}")

    print(f"Webhooks injected: {webhook_count}")
    print(f"v1beta1 schemas patched: {schema_count}")


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <crds-directory>", file=sys.stderr)
        sys.exit(2)
    main(sys.argv[1])
