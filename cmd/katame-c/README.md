# libkatame — C-shared library

This directory builds Katame into a native shared library
(`.so` / `.dylib` / `.dll`) usable from C, Python (ctypes), Ruby (FFI),
Node.js (node-ffi-napi), Java (JNA), and any other language with a C FFI.

The exported C ABI is intentionally narrow: every entry point takes a
single null-terminated UTF-8 JSON string and returns a single
`malloc`-allocated UTF-8 JSON string. Callers must release each returned
string with `kt_free`.

## Building

```bash
# Linux:
go build -buildmode=c-shared -tags cshared -o libkatame.so ./cmd/katame-c

# macOS:
go build -buildmode=c-shared -tags cshared -o libkatame.dylib ./cmd/katame-c

# Windows (requires MinGW or MSVC):
go build -buildmode=c-shared -tags cshared -o libkatame.dll ./cmd/katame-c
```

The build tag `cshared` is required — it excludes this `package main` from
the default `go build ./...` so the rest of the repository continues to
build without cgo or a C toolchain.

The build produces two artifacts:

- `libkatame.{so,dylib,dll}` — the shared library.
- `libkatame.h`              — the auto-generated C header with the
  exported function signatures.

## Exported functions

| Symbol        | Input                | Returns                              |
|---------------|----------------------|--------------------------------------|
| `kt_version`  | (none)               | semver string                        |
| `kt_protect`  | JSON options object  | JSON `{ok, bundle_path, manifest, generated_key_b64?}` |
| `kt_decrypt`  | JSON options object  | JSON `{ok, manifest, plaintext_b64?}` |
| `kt_inspect`  | bundle path (string) | JSON `{ok, manifest}`                |
| `kt_sign`     | JSON options object  | JSON `{ok, algorithm, signed_at, signature_b64}` |
| `kt_verify`   | JSON options object  | JSON `{ok, valid, algorithm, signed_at, manifest}` |
| `kt_free`     | pointer to free      | (void)                               |

On error, the JSON returned by every entry point looks like
`{"ok":false,"error":"<message>"}`.

See the package doc comment in `main.go` for the precise JSON shape of each
options object.

## Example (C)

```c
#include <stdio.h>
#include "libkatame.h"

int main(void) {
    char* json = kt_version();
    printf("katame version: %s\n", json);
    kt_free(json);
    return 0;
}
```

Link with `-L. -lkatame` (Linux/macOS) or with `katame.lib` on Windows.
