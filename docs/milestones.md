# Milestones

VaultPack grew as numbered milestones. Markers appear in CLI registration, package comments, and manifest fields.

## Completed (tagged in code)

| ID | Feature | Status |
| -- | ------- | ------ |
| **M16** | Config file + profiles | Done |
| **M19** | Plugin system | Done |
| **M22** | Key rotation & rewrap | Done |
| **M23** | Policy / RBAC | Done (CLI now **fail-closed** on load error) |
| **M24** | Sigstore Rekor transparency | Done |
| **M26** | Service mode (`vaultpack serve`) | Done (`/metrics` auth by default; `--metrics-public` opt-out) |
| **M27** | Tamper-evident repo | Done |
| **M28** | Key-source URIs | Done — portable + **dpapi://** (Windows), **keychain://** (macOS), **pkcs11://** / **piv://** file shims |
| **M30** | Shell completions | Done |

## Earlier layers (inferred, unnumbered historically)

Core CLI, streaming AEAD, integrity/audit, PQ crypto, cloud+KMS, batch/parallel, SDK/bindings, doctor/menu, hash-chained audit/`report`.

## Just shipped (this pass)

- **M28 complete:** platform key sources (`dpapi://`, `keychain://`, `pkcs11://` / `piv://` shims) — see `internal/keysource`
- CLI policy **fail-closed** when a policy path is configured
- Serve `/metrics` requires auth unless `--metrics-public`
- C ABI / Python / WASM parity for hybrid, KMS, compress, Shamir, `vp_rewrap`
- Dockerfile Go builder image aligned to **1.25**

## Backlog (1000+)

See **[milestones_backlog.md](milestones_backlog.md)** for **M31–M1041** concrete improvement items from a full-repo code review (crypto, CLI, serve, policy, audit, repo, cloud, KMS, plugins, SDK/bindings, CI, security, UX, compliance, packaging, …).

When completing an item, note it in the “Completed” table above and mark it done in the backlog.
