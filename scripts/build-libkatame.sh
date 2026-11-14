#!/usr/bin/env bash
# Build libkatame c-shared artifacts for the current OS/arch (M230).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="${1:-$ROOT/dist/libkatame}"
mkdir -p "$OUT"

GOOS="$(go env GOOS)"
GOARCH="$(go env GOARCH)"
case "$GOOS" in
  linux)  EXT=so ;;
  darwin) EXT=dylib ;;
  windows) EXT=dll ;;
  *) echo "unsupported GOOS=$GOOS"; exit 1 ;;
esac

LIB="libkatame.${EXT}"
HDR="libkatame.h"
echo "building $LIB ($GOOS/$GOARCH)…"
CGO_ENABLED=1 go build -buildmode=c-shared -tags cshared \
  -trimpath -ldflags "-s -w" \
  -o "$OUT/$LIB" "$ROOT/cmd/katame-c"

# go build -buildmode=c-shared writes the header next to the .so with the same basename.
if [[ -f "$OUT/libkatame.h" ]]; then
  :
elif [[ -f "$OUT/${LIB%.${EXT}}.h" ]]; then
  mv "$OUT/${LIB%.${EXT}}.h" "$OUT/libkatame.h"
fi

ARCHIVE="libkatame-${GOOS}-${GOARCH}"
(
  cd "$OUT"
  if [[ "$GOOS" == windows ]]; then
    zip -q "${ARCHIVE}.zip" "$LIB" libkatame.h
    echo "wrote $OUT/${ARCHIVE}.zip"
  else
    tar -czf "${ARCHIVE}.tar.gz" "$LIB" libkatame.h
    echo "wrote $OUT/${ARCHIVE}.tar.gz"
  fi
)
