#!/usr/bin/env bash
# Build katame.wasm + wasm_exec.js into bindings/typescript (M233).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DEST="$ROOT/bindings/typescript"
mkdir -p "$DEST"

echo "building katame.wasm…"
GOOS=js GOARCH=wasm go build -trimpath -ldflags "-s -w" \
  -o "$DEST/katame.wasm" "$ROOT/cmd/katame-wasm"

WASM_EXEC="$(go env GOROOT)/lib/wasm/wasm_exec.js"
if [[ ! -f "$WASM_EXEC" ]]; then
  WASM_EXEC="$(go env GOROOT)/misc/wasm/wasm_exec.js"
fi
cp "$WASM_EXEC" "$DEST/wasm_exec.js"
ls -la "$DEST/katame.wasm" "$DEST/wasm_exec.js"
