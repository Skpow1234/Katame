# Key-source URIs

Load symmetric and signing keys from URI-style sources. This is the stable abstraction for future HSM / keychain backends.

```bash
export VPACK_DATA_KEY="b64:$(base64 -w0 secret.key.raw)"
vaultpack protect --in secret.txt --out secret.vpack --key-source env://VPACK_DATA_KEY
vaultpack decrypt --in secret.vpack --out secret.out --key-source env://VPACK_DATA_KEY

vaultpack decrypt --in secret.vpack --out secret.out --key-source file://secret.key

export VPACK_SIGNING_KEY="$(cat signing.key)"
vaultpack sign --in secret.vpack --signing-key-source env://VPACK_SIGNING_KEY
```

## Supported today

| Scheme | Purpose |
| ------ | ------- |
| `file://` | Key file or PEM private key |
| `env://` | Environment variable |
| `b64://` | Inline base64 key material |

## Reserved (not implemented)

These return a clear “not available in this build” error:

| Scheme | Planned backend |
| ------ | --------------- |
| `pkcs11://` | PKCS#11 HSM |
| `piv://` | PIV / YubiKey |
| `keychain://` | macOS Keychain / Secret Service |
| `dpapi://` | Windows DPAPI |

`vaultpack doctor` reports which schemes are available.
