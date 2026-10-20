# SDK & language bindings

VaultPack ships libraries for Go, Python, and JavaScript so you can protect, decrypt, sign, verify, and inspect `.vpack` bundles without spawning the CLI.

## Go SDK (`pkg/vaultpack`)

Semver-stable (`vaultpack.Version`). Produces byte-identical bundles to the CLI for matching inputs.

```go
import "github.com/Skpow1234/Vaultpack/pkg/vaultpack"

res, err := vaultpack.Protect(vaultpack.ProtectOptions{
    InputPath:  "secret.txt",
    OutputPath: "secret.vpack",
})
os.WriteFile("secret.key", res.GeneratedKey, 0o600)

dec, err := vaultpack.Decrypt(vaultpack.DecryptOptions{
    InputPath:  "secret.vpack",
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

## Native shared library (C-shared)

```bash
go build -buildmode=c-shared -tags cshared -o libvaultpack.so ./cmd/vaultpack-c
```

Exports: `vp_version`, `vp_protect`, `vp_decrypt`, `vp_inspect`, `vp_sign`, `vp_verify`, `vp_free`. See [`cmd/vaultpack-c/README.md`](../cmd/vaultpack-c/README.md).

## Python (`bindings/python`)

```bash
pip install vaultpack
export VAULTPACK_LIB=/path/to/libvaultpack.so
```

```python
import base64, vaultpack
res = vaultpack.protect(input_path="secret.txt", output_path="secret.vpack")
key = base64.b64decode(res["generated_key_b64"])
out = vaultpack.decrypt(input_path="secret.vpack", key=key)
```

## TypeScript / WASM (`bindings/typescript`)

```bash
GOOS=js GOARCH=wasm go build -o vaultpack.wasm ./cmd/vaultpack-wasm
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" bindings/typescript/
```

```javascript
import { loadVaultpack } from "@vaultpack/wasm";
const Vaultpack = await loadVaultpack(new URL("./vaultpack.wasm", import.meta.url));
const enc = Vaultpack.protect({ plaintext: new TextEncoder().encode("hi") });
const dec = Vaultpack.decrypt({ bundle: enc.bundle, key: enc.generatedKey });
```

## Semver policy

Public Go SDK + on-disk `.vpack` format follow SemVer 2.0. Compatible additions are minor; breaking changes require a major bump. `internal/...` is not covered.
