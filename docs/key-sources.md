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

## Platform / HSM schemes (M28, M251–M252)

| Scheme | Platform | Behavior |
| ------ | -------- | -------- |
| `dpapi://path` or `dpapi://b64:…` | **Windows** | DPAPI unprotect → Katame key / PEM |
| `keychain://SERVICE/ACCOUNT` | **macOS** | `security find-generic-password -w` |
| `pkcs11://label[?pin=&slot=]` | All | File shim via `$KATAME_PKCS11_DIR`; native SoftHSM when built with `-tags katame_pkcs11` + `$KATAME_PKCS11_MODULE` |
| `piv://slot[?pin=]` | All | File shim via `$KATAME_PIV_DIR`; native YubiKey when built with `-tags katame_piv` |

`katame doctor` reports which schemes are available for this binary.

### Native SoftHSM (M251)

```bash
# Linux/macOS with SoftHSM2 + CGO
CGO_ENABLED=1 go build -tags katame_pkcs11 -o bin/katame ./cmd/katame
export KATAME_PKCS11_MODULE=/usr/lib/softhsm/libsofthsm2.so
export KATAME_PKCS11_PIN=1234
katame decrypt --in secret.katame --out secret.out --key-source 'pkcs11://my-aes-label'
katame sign --in secret.katame --signing-key-source 'pkcs11://my-ec-label?pin=1234'
```

Symmetric keys must be extractable AES-256 objects. Signing uses a PKCS#11 `crypto.Signer` (ECDSA/RSA-PSS) without exporting the private key.

### Native YubiKey PIV (M252)

```bash
CGO_ENABLED=1 go build -tags katame_piv -o bin/katame ./cmd/katame
export KATAME_PIV_PIN=123456
katame sign --in secret.katame --signing-key-source 'piv://9c'
```

Slots: `9a` (auth), `9c` (signature), `9d` (key management), `9e` (card auth). Private key material never leaves the device; CLI uses `keysource.ResolveSigner`.

### Windows DPAPI helper

Use `keysource.ProtectDPAPI` (Go) to create blobs, then store them and reference with `dpapi://path/to/blob`.
