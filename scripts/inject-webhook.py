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

def main(crd_dir: str) -> None:
    crd_path = pathlib.Path(crd_dir)
    if not crd_path.is_dir():
        print(f"ERROR: {crd_dir} is not a directory", file=sys.stderr)
        sys.exit(1)

    patched = 0
    for f in sorted(crd_path.glob("*.yaml")):
        try:
            doc = yaml.safe_load(f.read_text())
        except Exception as e:
            print(f"  SKIP  {f.name}: cannot parse YAML ({e})")
            continue

        if doc is None or doc.get("kind") != "CustomResourceDefinition":
            continue

        spec = doc.setdefault("spec", {})
        spec["conversion"] = WEBHOOK
        f.write_text(yaml.dump(doc, default_flow_style=False, sort_keys=False))
        patched += 1
        print(f"  webhook → {f.name}")

    print(f"Patched {patched} CRD(s)")


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <crds-directory>", file=sys.stderr)
        sys.exit(2)
    main(sys.argv[1])
