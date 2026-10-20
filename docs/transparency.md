# Transparency & Sigstore

VaultPack can upload bundle signatures to [Sigstore Rekor](https://docs.sigstore.dev/logging/overview/) and embed inclusion metadata in the manifest.

## Flags

| Flag | Description |
| ---- | ----------- |
| `--transparency` | Upload to Rekor |
| `--rekor-url` | Default `https://rekor.sigstore.dev` |
| `--keyless` | Fulcio keyless signing (requires `--transparency`) |
| `--fulcio-url` | Default `https://fulcio.sigstore.dev` |
| `--oidc-token-file` | OIDC ID token (or `VAULTPACK_OIDC_TOKEN` / `SIGSTORE_ID_TOKEN`) |
| `--check-transparency` | Verify Rekor inclusion + SET on `verify` |
| `--rekor-pubkey` | Offline Rekor public key (PEM) |

## Sign with personal key + Rekor

```bash
vaultpack sign --in artifact.vpack --signing-priv id.key --transparency
```

## Keyless via Fulcio

```bash
export SIGSTORE_ID_TOKEN=$(cat my-oidc.jwt)

vaultpack sign --in artifact.vpack \
  --keyless \
  --transparency
```

Keyless generates an ephemeral ECDSA P-256 key, obtains a short-lived Fulcio cert, signs, and uploads to Rekor.

## Verify

```bash
vaultpack verify --in artifact.vpack --pubkey id.pub --check-transparency
vaultpack verify --in artifact.vpack --pubkey id.pub --check-transparency --rekor-pubkey rekor.pub
```

Checks:

1. Detached signature valid against `--pubkey`
2. For each `transparency[]` entry: data hash, signature bytes, and SET against the log key

## Algorithm support

Rekor uses stdlib `crypto/x509`. `--transparency` is limited to: `ed25519`, `ecdsa-p256`, `ecdsa-p384`, `rsa-pss-2048`, `rsa-pss-4096`. Post-quantum / Ed448 are rejected with a clear error.

`attest --transparency` uploads a hashedrekord for provenance; the Rekor UUID is printed but not embedded in the manifest (use `sign --transparency` for that).
