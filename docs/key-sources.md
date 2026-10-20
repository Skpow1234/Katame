# Key-source URIs

Load symmetric and signing keys from URI-style sources (`--key-source`, `--signing-key-source`).

```bash
export KATAME_DATA_KEY="b64:$(base64 -w0 secret.key.raw)"
katame protect --in secret.txt --out secret.katame --key-source env://KATAME_DATA_KEY
katame decrypt --in secret.katame --out secret.out --key-source env://KATAME_DATA_KEY

katame decrypt --in secret.katame --out secret.out --key-source file://secret.key

export KATAME_SIGNING_KEY="$(cat signing.key)"
katame sign --in secret.katame --signing-key-source env://KATAME_SIGNING_KEY
```

## Portable schemes

| Scheme | Purpose |
| ------ | ------- |
| `file://` | Key file or PEM private key |
| `env://` | Environment variable |
| `b64://` | Inline base64 key material |

## Platform / HSM schemes (M28)

| Scheme | Platform | Behavior |
| ------ | -------- | -------- |
| `dpapi://path` or `dpapi://b64:…` | **Windows** | DPAPI unprotect → Katame key / PEM |
| `keychain://SERVICE/ACCOUNT` | **macOS** | `security find-generic-password -w` |
| `pkcs11://label` | All (shim) | Reads `$KATAME_PKCS11_DIR/<label>.key` or `.pem` |
| `piv://slot` | All (shim) | Reads `$KATAME_PIV_DIR/<slot>.pem` (signing); not for symmetric DEKs |

`katame doctor` reports which schemes are `on`/`off` for this binary.

Native SoftHSM / YubiKey PIV CGO bindings remain on the [backlog](milestones_backlog.md); the URI surface and file-backed shims are stable for CI and integration today.

### Windows DPAPI helper

Use `keysource.ProtectDPAPI` (Go) to create blobs, then store them and reference with `dpapi://path/to/blob`.
