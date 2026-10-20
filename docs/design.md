# Bundle design

## On-disk layout

A `.vpack` file is a ZIP archive:

```text
artifact.vpack
├── payload.bin        # ciphertext (chunked AEAD)
├── manifest.json      # encryption params, hashes, metadata
├── signature.sig      # optional detached signature
└── provenance.json    # optional SLSA-style provenance (attest --embed)
```

## Manifest versions

- **v1** — basic bundles
- **v2** — compression, multi-recipient, key splitting (and later fields)

Readers accept both. Current writer uses v2 when features require it.

## Manifest fields (overview)

Top-level (`internal/bundle.Manifest`):

| Field | Description |
| ----- | ----------- |
| `version` | `v1` / `v2` |
| `created_at` | RFC 3339 |
| `input` | Original name + size |
| `plaintext_hash` | Algo + digest |
| `encryption` | AEAD, nonce, tag, AAD, key ID, chunk size, KDF, hybrid, KMS wrap |
| `ciphertext` | Size / layout of `payload.bin` |
| `signature_algo` / `signed_at` | Optional signature metadata |
| `compression` | Optional gzip/zstd |
| `key_splitting` | Optional Shamir metadata |
| `rotated_from[]` | Rotation lineage (M22) |
| `transparency[]` | Rekor inclusion proofs (M24) |

Canonical JSON of the manifest (plus payload digest) is what signatures cover. See `internal/bundle` for exact struct tags and validation rules.

## Public SDK types

`pkg/vaultpack` re-exports these types so CLI `--json` and SDK JSON share the same field names.
