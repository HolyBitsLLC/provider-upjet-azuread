#!/usr/bin/env python3
"""Inject conversion webhook into CRD YAMLs at build time.

Without spec.conversion.strategy: Webhook in the CRD YAML, Crossplane's
APIEstablisher skips CA bundle injection — the CRD stays at strategy: None
and the singleton-list conversion webhook never activates.

This script mirrors the behavior of `make kustomize-crds` which applies
package/kustomize/webhook.yaml via kustomize to inject the webhook config.
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

# CRDs whose v1beta1 version should be un-served to prevent SSA crashes.
# SSA's field manager references the v1beta1 CRD schema (arrays for
# singleton-list fields like web, optionalClaims, api) against v1beta2
# stored data (objects), causing "expected list, got map" errors.
# Un-serving v1beta1 prevents new v1beta1 submits while the version
# stays listed so stored v1beta1 objects can still be read/converted.
UNSERVE_V1BETA1 = {
    "applications.azuread.upbound.io",
    "conditionalaccess.azuread.upbound.io",
    "groups.azuread.upbound.io",
    "invitations.azuread.upbound.io",
    "serviceprincipals.azuread.upbound.io",
}


def main(crd_dir: str) -> None:
    crd_path = pathlib.Path(crd_dir)
    if not crd_path.is_dir():
        print(f"ERROR: {crd_dir} is not a directory", file=sys.stderr)
        sys.exit(1)

    webhook_count = 0
    unserved_count = 0
    for f in sorted(crd_path.glob("*.yaml")):
        try:
            doc = yaml.safe_load(f.read_text())
        except Exception as e:
            print(f"  SKIP  {f.name}: cannot parse YAML ({e})")
            continue

        if doc is None or doc.get("kind") != "CustomResourceDefinition":
            continue

        changed = False

        spec = doc.setdefault("spec", {})
        # Inject conversion webhook
        spec["conversion"] = WEBHOOK
        webhook_count += 1
        changed = True

        # Un-serve v1beta1 for singleton-list CRDs so SSA never sees the
        # v1beta1 array schema.
        group = spec.get("group", "")
        versions = spec.get("versions", [])
        if group in UNSERVE_V1BETA1 and len(versions) >= 2:
            for v in versions:
                if v["name"] == "v1beta1" and v.get("served", True):
                    v["served"] = False
                    unserved_count += 1
                    changed = True
                    break

        if changed:
            f.write_text(yaml.dump(doc, default_flow_style=False, sort_keys=False))

        flags = []
        if group in UNSERVE_V1BETA1 and any(v["name"] == "v1beta1" and not v.get("served", True) for v in versions):
            flags.append("unserved-v1beta1")
        flags.append("webhook")
        print(f"  {','.join(flags)} → {f.name}")

    print(f"Webhooks injected: {webhook_count}")
    print(f"v1beta1 un-served: {unserved_count}")


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <crds-directory>", file=sys.stderr)
        sys.exit(2)
    main(sys.argv[1])
