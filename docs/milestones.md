# Milestones

VaultPack grew as numbered milestones. Markers appear in CLI registration, package comments, and manifest fields. There is no separate `ROADMAP.md`; this page reconstructs the map from the codebase.

## Explicitly tagged in code

| ID | Feature | Evidence |
| -- | ------- | -------- |
| **M16** | Config file + profiles (`--config`, `--profile`) | `internal/cli/root.go` |
| **M19** | Plugin system (subprocess KEM/sign) | `examples/plugin-dummy-kem`, `internal/plugin` |
| **M22** | Key rotation & rewrap (`rewrap`, `rotate-key`, add/remove recipient); `rotated_from` | `internal/cli/root.go`, `bundle.RotationEntry` |
| **M23** | Policy / RBAC (YAML/JSON/Rego) | `internal/cli/root.go`, `internal/policy` |
| **M24** | Sigstore Rekor transparency entries on manifests | `bundle.TransparencyEntry` comment |
| **M26** | Service mode (`vaultpack serve`) | `internal/cli/root.go`, `internal/serve` |
| **M27** | Tamper-evident repo (`vaultpack repo`) | `internal/cli/root.go`, `internal/repo` |
| **M30** | Shell completions | `internal/cli/root.go` |

## Earlier capability layers (inferred)

Not all early work is labeled `M*`, but the tree and commit history reflect roughly this order:

| Area | What shipped |
| ---- | ------------ |
| Core CLI | `protect` / `decrypt` / `inspect` / `hash` / `keygen` / `sign` / `verify` |
| Streaming crypto | Chunked AEAD, multiple ciphers, KDFs, Shamir |
| Integrity | `verify-integrity`, audit log, attest, Merkle seal |
| PQ / advanced crypto | ML-KEM, ML-DSA, SLH-DSA, Ed448, AES-GCM-SIV |
| Cloud + KMS | AWS/GCP/Azure KMS wrap; `az://` `s3://` `gs://` HTTPS I/O |
| Diff / batch / parallel | `diff`, batch-*, parallel chunk workers |
| SDK & bindings | `pkg/vaultpack`, C-shared, Python, WASM |
| UX | `doctor`, `menu`, key-source URIs, hash-chained audit + `report` |

## Likely open / reserved work

| Item | Status |
| ---- | ------ |
| HSM key sources (`pkcs11://`, `piv://`, `keychain://`, `dpapi://`) | Reserved URI schemes; not implemented |
| CLI policy fail-closed when policy path is set | Serve already fails closed; CLI still fail-open on load error |
| Docs tree | Previously referenced but missing; restored under `docs/` |

If you add a new milestone, tag it in `internal/cli/root.go` (or the owning package) and append a row here.
