# Tamper-evident repo

`vaultpack repo` maintains an append-only Merkle history of `.vpack` bundles. Each `repo add` appends a leaf, records a signed root, and can optionally anchor that root in Rekor.

```bash
vaultpack keygen --out repo-signing
vaultpack repo init --dir ./vp-repo --signing-key repo-signing.key --store-bundles
vaultpack repo add --dir ./vp-repo --bundle secret.vpack --signing-key repo-signing.key --copy
vaultpack repo list --dir ./vp-repo
vaultpack repo verify --dir ./vp-repo --verify-key repo-signing.pub
vaultpack repo anchor --dir ./vp-repo --signing-key repo-signing.key --rekor-url https://rekor.sigstore.dev
```

## Layout

```text
vp-repo/
├── repo.json        # schema, created_at, signing-key fingerprint
├── entries.jsonl    # append-only bundle entries
├── roots.jsonl      # append-only signed Merkle roots
└── bundles/         # optional stored copies
```

## Merkle construction (RFC 6962-style)

- `leaf = SHA-256(0x00 || canonical-entry)`
- `node = SHA-256(0x01 || left || right)`

`repo verify` replays `entries.jsonl`, recomputes every root in `roots.jsonl`, checks the `prev_root` chain, and verifies signatures when `--verify-key` is supplied.
