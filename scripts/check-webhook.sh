#!/usr/bin/env bash
# Validate that the CRD conversion webhook is enabled in the Crossplane package
# metadata whenever singleton-list conversions are configured.
#
# Root cause of the July 2026 incident: package/crossplane.yaml was missing
# crossplane.webhook.conversion: true, causing CRD conversion strategy to be
# None. v1beta1 singleton-list arrays were stored directly as v1beta2 without
# conversion, producing empty [{}] shells for all nested fields.
set -euo pipefail

SL_FILE="config/old-singleton-list-apis.txt"
PKG_FILE="package/crossplane.yaml"

if [ ! -f "$SL_FILE" ]; then
  echo "PASS: no singleton-list conversions configured ($SL_FILE not found)"
  exit 0
fi

SL_COUNT=$(grep -cve '^\s*#' -e '^\s*$' "$SL_FILE" 2>/dev/null || echo 0)

if [ "$SL_COUNT" -eq 0 ]; then
  echo "PASS: no singleton-list conversions configured"
  exit 0
fi

echo "Found $SL_COUNT resource(s) with singleton-list conversions in $SL_FILE"

if [ ! -f "$PKG_FILE" ]; then
  echo "FAIL: $PKG_FILE not found"
  exit 1
fi

if ! grep -qx '\s*conversion:\s*true' "$PKG_FILE"; then
  echo ""
  echo "FAIL: $SL_FILE lists $SL_COUNT resource(s) with singleton-list conversions,"
  echo "but $PKG_FILE is missing:"
  echo ""
  echo "  crossplane:"
  echo "    webhook:"
  echo "      conversion: true"
  echo ""
  echo "Without this, Crossplane sets CRD conversion strategy to None and the"
  echo "singleton-list->embedded-object conversion never runs at runtime."
  echo "This results in empty [{}] stubs for all nested fields (web, optionalClaims, api)."
  exit 1
fi

echo "PASS: webhook conversion enabled for $SL_COUNT singleton-list resource(s)"
