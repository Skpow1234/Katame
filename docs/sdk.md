# SDK & language bindings

Katame ships libraries for Go, Python, and JavaScript so you can protect, decrypt, sign, verify, and inspect `.katame` bundles without spawning the CLI.

## Go SDK (`pkg/katame`)

Semver-stable (`katame.Version`). Produces byte-identical bundles to the CLI for matching inputs.

```go
import "github.com/Skpow1234/Katame/pkg/katame"

res, err := katame.Protect(katame.ProtectOptions{
    InputPath:  "secret.txt",
    OutputPath: "secret.katame",
})
os.WriteFile("secret.key", res.GeneratedKey, 0o600)

dec, err := katame.Decrypt(katame.DecryptOptions{
    InputPath:  "secret.katame",
    OutputPath: "secret.out",
    Key:        keyBytes,
})
```

| Function | Purpose |
| -------- | ------- |
| `Protect` | Encrypt file or bytes → bundle |
| `Decrypt` | Decrypt → file or bytes |
| `Inspect` / `InspectBytes` | Manifest only |
| `SignBundle` / `SignBytes` | Detached signature |
| `Verify` / `VerifyBytes` | Verify signature |
| `Rewrap` | KMS DEK rewrap (where exposed) |
| `SplitShares` / `CombineShares` | Shamir helpers |

**Supported modes (v1.0):** symmetric key, password KDFs, hybrid/multi-recipient, KMS-wrapped DEKs, compression, Shamir splitting, detached signatures (including PQ), Rekor helpers, provenance embedding.

C / Python / WASM bindings now accept the same hybrid, KMS, compress, and Shamir fields (plus `kt_rewrap` / `rewrap()`).

## Native shared library (C-shared)

```bash
go build -buildmode=c-shared -tags cshared -o libkatame.so ./cmd/katame-c
```

Exports: `kt_version`, `kt_protect`, `kt_decrypt`, `kt_inspect`, `kt_sign`, `kt_verify`, `kt_free`. See [`cmd/katame-c/README.md`](../cmd/katame-c/README.md).

## Python (`bindings/python`)

```bash
pip install katame
export KATAME_LIB=/path/to/libkatame.so
```

```python
import base64, katame
res = katame.protect(input_path="secret.txt", output_path="secret.katame")
key = base64.b64decode(res["generated_key_b64"])
out = katame.decrypt(input_path="secret.katame", key=key)
```

## TypeScript / WASM (`bindings/typescript`)

```bash
GOOS=js GOARCH=wasm go build -o katame.wasm ./cmd/katame-wasm
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" bindings/typescript/
```

```javascript
import { loadKatame } from "@katame/wasm";
const Katame = await loadKatame(new URL("./katame.wasm", import.meta.url));
const enc = Katame.protect({ plaintext: new TextEncoder().encode("hi") });
const dec = Katame.decrypt({ bundle: enc.bundle, key: enc.generatedKey });
```

## Semver policy

Public Go SDK + on-disk `.katame` format follow SemVer 2.0. Compatible additions are minor; breaking changes require a major bump. `internal/...` is not covered.
