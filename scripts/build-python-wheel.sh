#!/usr/bin/env bash
# Build a platform wheel that vendors libkatame (M231).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="${1:-$ROOT/dist/wheels}"
LIBDIR="${2:-$ROOT/dist/libkatame}"
mkdir -p "$OUT"

GOOS="$(go env GOOS)"
case "$GOOS" in
  linux)  LIB=libkatame.so ;;
  darwin) LIB=libkatame.dylib ;;
  windows) LIB=libkatame.dll ;;
  *) echo "unsupported GOOS=$GOOS"; exit 1 ;;
esac

if [[ ! -f "$LIBDIR/$LIB" ]]; then
  bash "$ROOT/scripts/build-libkatame.sh" "$LIBDIR"
fi

PKG="$ROOT/bindings/python/katame"
cp "$LIBDIR/$LIB" "$PKG/$LIB"
cleanup() { rm -f "$PKG/$LIB"; }
trap cleanup EXIT

python3 -m pip install --upgrade build wheel >/dev/null
(
  cd "$ROOT/bindings/python"
  python3 -m build --wheel --outdir "$OUT"
)
echo "wheels in $OUT"
ls -la "$OUT"
