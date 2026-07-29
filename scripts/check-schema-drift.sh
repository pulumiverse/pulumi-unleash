#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

SCHEMA=provider/cmd/pulumi-resource-unleash/schema.json
BASELINE=schema-baseline.json

# .make/schema is a sentinel; without removing it `make tfgen` no-ops and the
# comparison below would pass without regenerating anything.
rm -f .make/schema
make tfgen >/dev/null

if ! diff -q "$BASELINE" "$SCHEMA" >/dev/null; then
  echo "Schema drift detected. Review, then update $BASELINE deliberately:"
  diff -u "$BASELINE" "$SCHEMA" | head -60
  exit 1
fi
echo "Schema matches baseline."
