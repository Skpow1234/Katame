# CLI usage

## Quick examples

```bash
katame protect --in config.json
katame protect --in config.json --cipher chacha20-poly1305
katame protect --in config.json --password "pass" --kdf scrypt
katame decrypt --in config.json.katame --out config.json --key config.json.key

katame keygen --out alice --algo x25519-aes-256-gcm
katame protect --in config.json --recipient alice.pub
katame decrypt --in config.json.katame --out config.json --privkey alice.key

katame protect --in config.json --kms-provider aws --kms-key-id alias/my-key
katame inspect --in config.json.katame
katame verify-integrity --in config.json.katame --key config.json.key
```

## Global flags

| Flag | Description |
| ---- | ----------- |
| `--json` | Output results as JSON |
| `--quiet` | Minimal output (errors only) |
| `--verbose` | Debug logging |
| `--audit-log` | Audit log file (or `KATAME_AUDIT_LOG`) |
| `--config` | Config file (or `KATAME_CONFIG`) |
| `--profile` | Profile: `dev`, `staging`, `prod` (or `KATAME_PROFILE`) |
| `--policy` | Policy file (or `KATAME_POLICY`) |
| `--version` | Print version |

Cloud flags (`--azure-account`, `--aws-region`, …) are documented in [cloud.md](cloud.md). Transparency flags are in [transparency.md](transparency.md).

## Configuration

Search order: `--config` / `KATAME_CONFIG` → `~/.katame.yaml` → `./.katame.yaml`.

Precedence: CLI flags > environment > config/profile > defaults.

```bash
katame config
katame config --json
katame --config ./katame.yaml --profile prod protect --in data.csv
```

YAML keys: `audit_log`, `cipher`, `chunk_size`, `output_dir`, `default_key_path`, `default_pubkey_path`, `recipients`, plus `profiles.<name>.*`.

## Commands

### `protect`

Encrypt a file into a `.katame` bundle.

| Flag | Default | Description |
| ---- | ------- | ----------- |
| `--in` | (required) | Input file |
| `--out` | `<input>.katame` | Output bundle |
| `--key-out` | `<input>.key` | Generated key path |
| `--key` | | Existing key (skip generation) |
| `--aad` | | Additional authenticated data |
| `--cipher` | `aes-256-gcm` | AEAD cipher |
| `--hash-algo` | `sha256` | Plaintext hash |
| `--sign` / `--signing-priv` / `--sign-algo` | | Sign in one step |
| `--password` / `--password-file` / `--kdf` | | Password mode |
| `--recipient` | | Hybrid recipient PEM (repeatable) |
| `--compress` | `none` | `none`, `gzip`, `zstd` |
| `--split-shares` / `--split-threshold` | | Shamir split |
| `--kms-provider` / `--kms-key-id` | | KMS wrap DEK |
| `--stdin` / `--stdout` | | Pipe mode |

Ciphers: `aes-256-gcm`, `chacha20-poly1305`, `xchacha20-poly1305`, `aes-256-gcm-siv` (all 32-byte keys, 64 KiB chunked streaming).

### `decrypt`

Provide exactly one of `--key`, `--password`, `--privkey`, or `--kms-provider` (when the bundle used KMS).

### `inspect` / `diff` / `hash`

- `inspect --in <bundle> [--json] [--redact]`
- `diff --a <b1> --b <b2> [--ignore-time] [--ignore-nonce] [--json]` — exit `0` identical, `10` differ
- `hash --in <file> [--algo sha256]` — also sha512, sha3-*, blake2b-*, blake3

### `keygen` / `sign` / `verify` / `verify-integrity`

Signing algos: Ed25519, Ed448, ECDSA P-256/P-384, RSA-PSS 2048/4096, ML-DSA-65/87, SLH-DSA-128s.

Encryption keygen: X25519, ECIES-P256, RSA-OAEP 2048/4096, ML-KEM-768/1024 (+ plugins).

`verify` / `verify-integrity` exit `0` on success, `10` on failure.

### Shamir: `split-key` / `combine-key`

```bash
katame split-key --in data.key --shares 5 --threshold 3
katame combine-key --share data.key.share1 --share data.key.share3 --share data.key.share5 --out recovered.key
```

### Batch

```bash
katame batch-protect --dir ./exports/ --out-dir ./encrypted/ --workers 4
katame batch-decrypt --dir ./encrypted/ --out-dir ./decrypted/ --key ./encrypted/batch.key
katame batch-inspect --dir ./encrypted/
```

### UX / ops

- `katame doctor` — runtime diagnostics
- `katame menu` — interactive TUI
- `katame completion bash|zsh|fish|powershell`

## Exit codes

| Code | Meaning |
| ---- | ------- |
| `0` | Success |
| `1` | Generic error |
| `2` | Invalid arguments |
| `10` | Verification failed |
| `11` | Decryption failed |
| `12` | Unsupported version or algorithm |

## Related

- [Cloud](cloud.md) · [Rotation](rotation.md) · [Audit](audit.md) · [Policy](policy.md)
- [Repo](repo.md) · [Serve](serve.md) · [Recipes](recipes.md)
