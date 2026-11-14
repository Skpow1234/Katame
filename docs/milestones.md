# Milestones

Katame grew as numbered milestones. Markers appear in CLI registration, package comments, and manifest fields.

## Completed (tagged in code)

| ID | Feature | Status |
| -- | ------- | ------ |
| **M16** | Config file + profiles | Done |
| **M19** | Plugin system | Done |
| **M22** | Key rotation & rewrap | Done |
| **M23** | Policy / RBAC | Done (CLI now **fail-closed** on load error) |
| **M24** | Sigstore Rekor transparency | Done |
| **M26** | Service mode (`katame serve`) | Done (`/metrics` auth by default; `--metrics-public` opt-out) |
| **M27** | Tamper-evident repo | Done |
| **M28** | Key-source URIs | Done — portable + **dpapi://** (Windows), **keychain://** (macOS), **pkcs11://** / **piv://** file shims |
| **M30** | Shell completions | Done |

## Earlier layers (inferred, unnumbered historically)

Core CLI, streaming AEAD, integrity/audit, PQ crypto, cloud+KMS, batch/parallel, SDK/bindings, doctor/menu, hash-chained audit/`report`.

## Just shipped (this pass)

- **M28 complete:** platform key sources (`dpapi://`, `keychain://`, `pkcs11://` / `piv://` shims) — see `internal/keysource`
- CLI policy **fail-closed** when a policy path is configured
- Serve `/metrics` requires auth unless `--metrics-public`
- C ABI / Python / WASM parity for hybrid, KMS, compress, Shamir, `kt_rewrap`
- **M31–M48:** crypto fuzz coverage (ChaCha/XChaCha/GCM-SIV/hybrid/PQ/SSS/KDF), aesgcmsiv `Open` errors, hard chunk-size cap, truncation/reorder tests, `Zeroize` / CT fingerprint helpers
- **M73–M80:** CLI commands return errors instead of `os.Exit`; `Execute` maps sentinels (`ExitPolicyDenied=13`, `ErrBundlesDiffer`→10)
- **M251–M252 / M256:** optional native SoftHSM (`-tags katame_pkcs11`) and YubiKey PIV (`-tags katame_piv`); `ResolveSigner`; doctor reports shim vs native details
- **M115–M132 (partial):** serve `/readyz`, `/v1/{keygen,hash,verify-integrity,attest}`, `X-Request-Id`, `error_code`, principal rate limits, protect concurrency cap, OpenAPI Rewrap + CI lint; deferred multipart/OTel/KMS ctx
- **M230–M233:** release ships `libkatame` archives, Python wheels with vendored lib + `.pyi`, WASM tarball / optional npm publish
- **M81–M98 / M103–M106 / M113–M114:** config→protect wiring, `policy_file`, dry-run, completions, `keygen --list-algos`, doctor probes, CLI tests
- **M99–M112:** decrypt `--stdin`, azureio→cloudio, encrypt progress, menu serve/repo/policy/doctor, `make man` / `cmd/gendocs`, Windows key ACL + tests, `RedactSecrets`, `convert`, batch policy, `--force`/`refuseIfExists`

## Backlog (1000+)

See **[milestones_backlog.md](milestones_backlog.md)** for **M31–M1041** concrete improvement items from a full-repo code review (crypto, CLI, serve, policy, audit, repo, cloud, KMS, plugins, SDK/bindings, CI, security, UX, compliance, packaging, …).

When completing an item, note it in the “Completed” table above and mark it done in the backlog.
