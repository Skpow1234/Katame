# Milestone backlog (M31+)

Generated from a full-repo improvement pass. Items are actionable;
priority is not implied by number order. Completed historical milestones:
see [milestones.md](milestones.md).

| ID | Milestone |
| -- | --------- |
| **M31** | Add fuzz for ChaCha20-Poly1305 decrypt/round-trip |
| **M32** | Add fuzz for XChaCha20-Poly1305 decrypt/round-trip |
| **M33** | Add fuzz for AES-256-GCM-SIV seal/open |
| **M34** | Replace panic with errors in aesgcmsiv size checks |
| **M35** | Add fuzz for hybrid encapsulate/decapsulate (X25519) |
| **M36** | Add fuzz for hybrid ECIES-P256 |
| **M37** | Add fuzz for hybrid RSA-OAEP |
| **M38** | Add fuzz for ML-KEM encapsulate/decapsulate |
| **M39** | Add fuzz for ML-DSA sign/verify |
| **M40** | Add fuzz for SLH-DSA sign/verify |
| **M41** | Add fuzz for Shamir corrupt-share rejection |
| **M42** | Add fuzz for KDF weak-param rejection |
| **M43** | Enforce MaxChunkSizeRecommended as hard cap |
| **M44** | Property-test last-chunk truncation for all ciphers |
| **M45** | Add reorder/duplicate chunk attack tests |
| **M46** | Race-test parallel stream encrypt/decrypt at high workers |
| **M47** | Zeroize DEK slices after AEAD use |
| **M48** | Constant-time key fingerprint compare helper |
| **M49** | Document nonce derivation for 12B vs 24B nonces |
| **M50** | Bench ML-KEM encapsulate |
| **M51** | Bench ML-DSA sign |
| **M52** | Negative tests for unsupported cipher/hash/KDF strings |
| **M53** | Validate SLH-DSA naming vs NIST in docs |
| **M54** | Document Ed448 is signing-only (no hybrid) |
| **M55** | Streaming compress+encrypt without full plaintext buffer |
| **M56** | Golden corrupt-tag corpus per AEAD |
| **M57** | Deprecate legacy Ed25519 raw key format with warn |
| **M58** | Add HPKE (RFC 9180) hybrid scheme |
| **M59** | AAD mismatch tests for every cipher including GCM-SIV |
| **M60** | Document Argon2id parallelism vs --parallel-workers |
| **M61** | Publish JSON Schema for manifest.json |
| **M62** | Manifest version-migration helper |
| **M63** | Fuzz ZIP zip-slip / unexpected entries |
| **M64** | Reject duplicate ZIP entries for manifest/payload |
| **M65** | Optional content-type on InputMeta |
| **M66** | Canonical JSON golden vectors across languages |
| **M67** | CLI --json inspect vs SDK Inspect field parity tests |
| **M68** | Document rotated_from invalid-proof rejection |
| **M69** | signature_algo nil-vs-ed25519 legacy decoder tests |
| **M70** | Optional encrypted metadata without full decrypt |
| **M71** | Validate provenance.json in bundle.Validate |
| **M72** | Format capability bits for forward-compat fail-closed |
| **M73** | Stop os.Exit in decrypt; return errors |
| **M74** | Stop os.Exit in verify; return errors |
| **M75** | Stop os.Exit in hash; return errors |
| **M76** | Stop os.Exit in verify-integrity; return errors |
| **M77** | Stop os.Exit in diff; return errors |
| **M78** | Stop os.Exit in repo CLI; return errors |
| **M79** | Stop os.Exit in audit export; return errors |
| **M80** | Dedicated exit code for ErrPolicyDenied |
| **M81** | Thread config Cipher into protect defaults |
| **M82** | Thread config ChunkSize into protect defaults |
| **M83** | Thread config Recipients into protect defaults |
| **M84** | Thread config KMS defaults into protect |
| **M85** | Add PolicyFile to config Profile struct |
| **M86** | Surface config.Load errors instead of ignoring |
| **M87** | CLI tests for serve command |
| **M88** | CLI tests for repo command |
| **M89** | CLI tests for menu command |
| **M90** | CLI tests for attest/seal |
| **M91** | CLI tests for add/remove recipient |
| **M92** | CLI tests for rewrap |
| **M93** | CLI tests for audit export |
| **M94** | CLI tests for completion |
| **M95** | Subprocess tests for exit codes |
| **M96** | Add --dry-run to protect |
| **M97** | Add --dry-run to rotate-key |
| **M98** | Add --dry-run to rewrap |
| **M99** | True stdin/stdout binary streaming without temp files |
| **M100** | Unify azureio.go with cloudio.go |
| **M101** | Progress reporting for large encrypt under --verbose |
| **M102** | Menu feature parity for serve/repo/policy/doctor |
| **M103** | Shell completion for cipher enums |
| **M104** | Shell completion for hash/KDF/algo enums |
| **M105** | keygen --list-algos including plugins |
| **M106** | Fix hash bad-algo exit code semantics |
| **M107** | Generate man pages from cobra |
| **M108** | Windows ACL key-file permission tests |
| **M109** | Redact PEM/password from verbose logs everywhere |
| **M110** | katame convert legacy Ed25519 to PKCS#8 |
| **M111** | Batch commands honor same policy ops |
| **M112** | --force vs refuse-overwrite consistency |
| **M113** | doctor probes plugin dir executables |
| **M114** | doctor probes KMS credential reachability |
| **M115** | OpenAPI schemas for Rewrap request/response |
| **M116** | Add /v1/keygen endpoint |
| **M117** | Add /v1/hash endpoint |
| **M118** | Add /v1/verify-integrity endpoint |
| **M119** | Add /v1/attest endpoint |
| **M120** | X-Request-Id on serve requests |
| **M121** | Structured error_code in JSON errors |
| **M122** | Multipart upload for large protect/decrypt |
| **M123** | GET /readyz separate from /healthz |
| **M124** | Pass request context into KMS unwrap |
| **M125** | mTLS+OIDC combination tests |
| **M126** | Per-principal (JWT sub) rate limits |
| **M127** | OpenTelemetry traces for protect/decrypt handlers |
| **M128** | Cap concurrent in-flight protect jobs |
| **M129** | UNIX socket umask tests |
| **M130** | Expand server_test for OIDC failure |
| **M131** | Expand server_test for rate-limit 429 |
| **M132** | Validate OpenAPI fixtures in CI |
| **M133** | CORS docs for WASM clients |
| **M134** | Audit logger reopen without restart |
| **M135** | HTTP 401 vs 403 for auth vs policy deny |
| **M136** | Policy schema v2 with default deny templates |
| **M137** | Policy conditions for compression algo |
| **M138** | Policy conditions for chunk size |
| **M139** | Policy conditions for hash algo |
| **M140** | Policy conditions for KMS provider name |
| **M141** | OIDC claims mapped into policy Context |
| **M142** | Fuzz Rego evaluation with malicious inputs |
| **M143** | Timezone midnight edge tests for TimeWindow |
| **M144** | Ship example deny-unsigned-decrypt policy |
| **M145** | Ship business-hours policy example |
| **M146** | Cache compiled Rego with mtime invalidate |
| **M147** | Remove unused sort import workaround in policy |
| **M148** | Propagate MatchedRule into serve audit |
| **M149** | Audit ops for repo-* commands |
| **M150** | Audit ops for serve-* commands |
| **M151** | Audit ops for diff/doctor/config/menu |
| **M152** | Complete OCSF class for every Op* |
| **M153** | Append-only hardening option for audit files |
| **M154** | Remote audit sink (HTTPS webhook) |
| **M155** | Syslog audit sink |
| **M156** | Merkle checkpoint signing for audit chain |
| **M157** | Concurrent FileLogger stress tests |
| **M158** | Export audit verify in pkg/katame |
| **M159** | Redact secrets from audit Error fields |
| **M160** | Encrypt-at-rest option for audit logs |
| **M161** | Report correlates audit with Rekor UUIDs |
| **M162** | Integration tests for repo rewrite.go |
| **M163** | Concurrent repo Add stress test |
| **M164** | Directory watch auto-add for .katame |
| **M165** | Repo prune/tombstone preserving history |
| **M166** | HTTP serve of repo inclusion proofs |
| **M167** | Persist partial Merkle tree on disk |
| **M168** | repo verify --offline with cached proofs |
| **M169** | Repo signing-key rotation without invalidating roots |
| **M170** | --require-signed-roots hard mode |
| **M171** | Fuzz Merkle inclusion proofs |
| **M172** | Content-addressed bundle store with hardlink dedup |
| **M173** | Unit tests for cloud/s3 |
| **M174** | Unit tests for cloud/gcs |
| **M175** | Unit tests for cloud/azureadp |
| **M176** | Streaming cloud Upload without full buffer |
| **M177** | Streaming cloud Download without full buffer |
| **M178** | S3 SSE-KMS option |
| **M179** | GCS CMEK option |
| **M180** | Azure customer-managed key option |
| **M181** | Retry with exponential backoff on cloud 5xx |
| **M182** | --allow-insecure-http gate for http:// |
| **M183** | Checksum verify after cloud download |
| **M184** | HTTPS prefix List or document unsupported |
| **M185** | doctor IRSA / workload-identity checks |
| **M186** | Multipart S3 upload for large outputs |
| **M187** | Unify URI parsing helper for Windows drives |
| **M188** | Cloud IT tests behind KATAME_IT_CLOUD |
| **M189** | Cache Store clients per scheme in batch |
| **M190** | Fix User-Agent katame/m21 string |
| **M191** | Context on KMS Provider interface |
| **M192** | Customizable AWS Encryption Context |
| **M193** | Customizable GCP KMS AAD |
| **M194** | Cache Azure Key Vault clients by URL |
| **M195** | Hashicorp Vault Transit KMS provider |
| **M196** | Oracle OCI KMS provider |
| **M197** | Soft-register cloud KMS only when creds exist |
| **M198** | LocalStack KMS integration tests |
| **M199** | Asymmetric cloud KMS wrap path |
| **M200** | Prometheus KMS latency/error metrics |
| **M201** | Register returns error on duplicate provider |
| **M202** | Document keyID format matrix per cloud |
| **M203** | Plugin package unit tests |
| **M204** | Stop DEK-via-env; use stdin for plugins |
| **M205** | Plugin process timeout + stderr capture |
| **M206** | Plugin binary hash allowlist |
| **M207** | Plugin protocol version negotiation |
| **M208** | Windows job-object sandbox guidance |
| **M209** | Unix seccomp sandbox guidance for plugins |
| **M210** | Multi-dir plugin discovery |
| **M211** | Surface plugins in doctor |
| **M212** | Example sign plugin |
| **M213** | Log skipped broken plugins at debug |
| **M214** | CI e2e for dummy KEM plugin |
| **M215** | Populate ProtectResult.BundleBytes from OutputWriter |
| **M216** | SDK rotate-key API |
| **M217** | SDK add-recipient API |
| **M218** | SDK remove-recipient API |
| **M219** | SDK VerifyIntegrity API |
| **M220** | SDK attest/provenance helpers |
| **M221** | SDK seal Merkle ops |
| **M222** | Standalone SplitShares export |
| **M223** | ProtectContext / DecryptContext variants |
| **M224** | Table tests for every cipher x KDF x compress |
| **M225** | Unit tests for internal/katameop |
| **M226** | Stable SDK sentinel errors |
| **M227** | Write API symmetric to Read for Bundle |
| **M228** | Experimental PQ opt-in gate |
| **M229** | examples/sdk-* programs |
| **M230** | Ship libkatame in GitHub Releases |
| **M231** | Python wheels vendoring shared lib |
| **M232** | Python .pyi type stubs |
| **M233** | Publish @katame/wasm to npm |
| **M234** | WASM size budget / trimmed build |
| **M235** | Playwright browser smoke for WASM |
| **M236** | Node ESM/CJS dual package tests |
| **M237** | pkg-config .pc for libkatame |
| **M238** | Ruby FFI example |
| **M239** | Java JNA example |
| **M240** | C# P/Invoke example |
| **M241** | Panic recovery in c-shared exports |
| **M242** | Fulcio client unit tests |
| **M243** | OIDC helper unit tests |
| **M244** | Rekor algo matrix tests for PQ uploads |
| **M245** | Cache Rekor public keys with expiry |
| **M246** | serve /v1/verify --check-transparency option |
| **M247** | Embed Fulcio cert chain in transparency entry |
| **M248** | Retry/backoff for Rekor 429 |
| **M249** | Private Rekor/Fulcio via config file |
| **M250** | Sigstore staging e2e (opt-in) |
| **M251** | Native SoftHSM PKCS#11 module binding |
| **M252** | Native YubiKey PIV via go-piv |
| **M253** | tpm:// sealed DEK support |
| **M254** | Windows file:// drive-letter round-trip CI |
| **M255** | age:// or SSH recipient key support |
| **M256** | doctor reports build tags for schemes |
| **M257** | Password memory lock / wipe after KDF |
| **M258** | Refuse --auth-disabled on non-loopback without confirm |
| **M259** | Sandbox ZIP max uncompressed size |
| **M260** | Crypto review checklist for Circl/OPA upgrades |
| **M261** | Threat model doc for plugin trust boundary |
| **M262** | --pin-hash expected digest on decrypt |
| **M263** | Overwrite-before-unlink temp cloud files |
| **M264** | SBOM attestation on Docker images |
| **M265** | Rebuild-from-SBOM verification docs |
| **M266** | Tee hash while streaming encrypt |
| **M267** | Adaptive chunk size by file size |
| **M268** | Tune parallel worker default vs GOMAXPROCS |
| **M269** | Buffer pools across EncryptStream chunks |
| **M270** | pprof endpoints behind auth on serve |
| **M271** | Cache parsed PEM pubs in multi-recipient |
| **M272** | Stream ZIP writing for bundles |
| **M273** | CI benchstat regression gate |
| **M274** | Standardize zerolog fields across CLI |
| **M275** | OpenTelemetry metrics parity with Prometheus |
| **M276** | Tracing for cloud download-decrypt-upload |
| **M277** | Correlate audit EntryHash with request ID |
| **M278** | VCS revision in /v1/version |
| **M279** | Structured /healthz detail payload |
| **M280** | Fuzz targets for hybrid/KDF/SSS/policy/ZIP |
| **M281** | Fix Makefile fuzz recursive package paths — **done** (Makefile + CI use single-package paths) |
| **M282** | Run golangci-lint in GitHub Actions — **done** (`.golangci.yml` + `ci.yml` lint job) |
| **M283** | Enable -race in CI test job — **done** (Unix runners; Windows remains non-race) |
| **M284** | Add linux/arm64 to CI build matrix — **done** |
| **M285** | Add windows/arm64 to CI build matrix — **done** |
| **M286** | Publish fuzz crashers as CI artifacts — **done** (upload on fuzz job failure) |
| **M287** | Coverage threshold for crypto/katameop |
| **M288** | Bindings smoke in CI |
| **M289** | Dependabot for Go modules and Actions — **done** (`.github/dependabot.yml`) |
| **M290** | Snap packaging |
| **M291** | Flatpak packaging |
| **M292** | macOS notarization |
| **M293** | Windows Authenticode signing |
| **M294** | Release notes auto-embed SBOM path |
| **M295** | Consumer cosign verify-blob CI example |
| **M296** | Cross-compile c-shared in release matrix |
| **M297** | Pin GitHub Actions by SHA |
| **M298** | Interactive first-run wizard |
| **M299** | Colored human diff for inspect |
| **M300** | Generate docs/usage.md from cobra in CI |
| **M301** | Troubleshooting runbook for KMS permissions |
| **M302** | Architecture diagram katameop/CLI/serve/bindings |
| **M303** | ADR folder for cipher defaults and policy |
| **M304** | docs/compat.md feature matrix |
| **M305** | Named pipe serve on Windows |
| **M306** | Long-path support for Windows batch trees |
| **M307** | macOS hardened runtime entitlements notes |
| **M308** | Warn on unknown config profile name |
| **M309** | katame config validate |
| **M310** | TOML/JSON config file support |
| **M311** | Unify KATAME_* vs KATAME_* env prefixes |
| **M312** | Secure temp file helper in util |
| **M313** | Replace panics in repo merkle subtree helpers |
| **M314** | Shared deterministic JSON encode helper |
| **M315** | Shared JSON envelope types pkg/katame/apienc |
| **M316** | Add /v1/diff endpoint |
| **M317** | Add /v1/batch-protect endpoint |
| **M318** | Add /v1/batch-decrypt endpoint |
| **M319** | Serve webhook for audit events |
| **M320** | Bundle format v3 planning document |
| **M321** | Optional AEAD key commitment field |
| **M322** | Age-compatible recipient stanza export |
| **M323** | Sigstore bundle format export alongside Rekor |
| **M324** | SLSA provenance v1.0 predicate upgrade |
| **M325** | in-toto attestation layout for attest |
| **M326** | SPIFFE identity for serve mTLS |
| **M327** | Workload identity federation doc for GKE/EKS/AKS |
| **M328** | HSM-backed repo root signing |
| **M329** | Offline air-gap mode package |
| **M330** | Deterministic protect mode for reproducible builds |
| **M331** | Content-defined chunking option |
| **M332** | Zstd dictionary training helper |
| **M333** | BLAKE3 keyed MAC option alongside AEAD |
| **M334** | Post-quantum hybrid combiner (X25519+ML-KEM) |
| **M335** | Algorithm agility policy deny for weak ciphers |
| **M336** | FIPS mode build tag restricting algos |
| **M337** | BoringCrypto / system crypto optional build |
| **M338** | Constant-time RSA-OAEP path audit |
| **M339** | Side-channel test notes for ChaCha |
| **M340** | Memory hard KDF auto-tune by host RAM |
| **M341** | Hardware AES-NI detection in doctor |
| **M342** | ARM crypto extension notes in docs |
| **M343** | WASM SIMD acceleration investigation |
| **M344** | TinyGo WASM experimental target |
| **M345** | Deno loader for katame.wasm |
| **M346** | Bun loader smoke test |
| **M347** | React Native / mobile binding feasibility |
| **M348** | Electron packaging example |
| **M349** | Grafana dashboard JSON for serve metrics |
| **M350** | Alerting rules for auth_denied_total |
| **M351** | SLO doc for serve latency histograms |
| **M352** | Chaos test: kill mid-encrypt recoverability |
| **M353** | Property-based tests with testing/quick |
| **M354** | Mutation testing for policy evaluator |
| **M355** | Contract tests CLI vs SDK vs serve vs C ABI |
| **M356** | Golden bundle corpus in testdata/ with licenses |
| **M357** | Backward-compat reader tests for every released tag |
| **M358** | Canary release channel documentation |
| **M359** | Feature flag file for experimental schemes |
| **M360** | Telemetry opt-in (anonymous usage counts) |
| **M361** | Crash reporter redacting secrets |
| **M362** | Local LLM-friendly docs/llms.txt index |
| **M363** | Spellcheck CI for docs/ |
| **M364** | Link checker CI for docs/ |
| **M365** | OpenAPI lint (spectral) in CI |
| **M366** | goreleaser nightlies on main |
| **M367** | SBOM diff in PRs |
| **M368** | License compliance scan (go-licenses) |
| **M369** | Pin Circl minor with renovate group |
| **M370** | OPA version bump playbook |
| **M371** | Azure SDK upgrade playbook |
| **M372** | AWS SDK upgrade playbook |
| **M373** | GCP SDK upgrade playbook |
| **M374** | Reduce charmbracelet TUI dep surface for non-menu builds |
| **M375** | Build tag nomainmenu to drop huh dependency |
| **M376** | Build tag nocloud to drop cloud SDKs |
| **M377** | Build tag nokms for minimal binary |
| **M378** | Build tag nopq for classical-only binary |
| **M379** | Size budget for linux-amd64 binary |
| **M380** | UPX prohibition / reproducible-build policy |
| **M381** | nix flake packaging |
| **M382** | Guix package recipe |
| **M383** | Chocolatey automatic updater tests |
| **M384** | Winget version schema validation |
| **M385** | Scoop autoupdate hash verification |
| **M386** | Homebrew bottle CI |
| **M387** | APT repo signing key rotation doc |
| **M388** | YUM repo metadata tests |
| **M389** | Alpine APKINDEX generation check |
| **M390** | Docker HEALTHCHECK for serve image |
| **M391** | distroless debug variant for troubleshooting |
| **M392** | Cosign keyless identity regex tests |
| **M393** | SLSA builder provenance for goreleaser |
| **M394** | Hermetic build with vendored modules option |
| **M395** | Offline module mirror docs |
| **M396** | Proxy/corporate CA docs for Rekor/Fulcio |
| **M397** | IPv6 listen tests for serve |
| **M398** | HTTP/2 cleartext h2c support decision |
| **M399** | gRPC gateway alternative to JSON HTTP |
| **M400** | Connect-RPC experiment |
| **M401** | WebSocket streaming protect experiment |
| **M402** | QUIC listen experiment |
| **M403** | Rate limit headers (Retry-After) |
| **M404** | W3C Trace Context propagation |
| **M405** | Baggage for tenant ID in serve |
| **M406** | Multi-tenant policy namespaces |
| **M407** | Per-tenant KMS key allowlists |
| **M408** | Admin API for hot-reload policy |
| **M409** | Admin API for hot-reload auth token |
| **M410** | Circuit breaker around cloud KMS |
| **M411** | Hedged requests for Rekor upload |
| **M412** | Idempotency-Key header for protect |
| **M413** | ETag for inspect responses |
| **M414** | Conditional rewrap if-match manifest hash |
| **M415** | Soft-delete recycle bin for repo bundles |
| **M416** | Repo GC dry-run |
| **M417** | Repo fsck command |
| **M418** | Repo export/import snapshot format |
| **M419** | Cross-repo Merkle bridge |
| **M420** | Threshold signing for repo roots (t of n) |
| **M421** | Multisig approve workflow for decrypt policy |
| **M422** | Break-glass policy override with dual control |
| **M423** | Time-lock encryption until timestamp |
| **M424** | Dead-man switch share release recipe |
| **M425** | Shamir share QR code export |
| **M426** | Paper key backup PDF generator |
| **M427** | NFC NDEF key share experiment |
| **M428** | YubiKey OTP as policy second factor |
| **M429** | WebAuthn gate for serve decrypt |
| **M430** | Passkey unlock for local menu |
| **M431** | TPM 2.0 PCR-bound keys |
| **M432** | Secure Enclave / Keychain ACL labels on Darwin |
| **M433** | Windows Hello / CNG key storage |
| **M434** | Linux kernel keyring key-source |
| **M435** | systemd credentials integration |
| **M436** | Kubernetes CSI secrets driver example |
| **M437** | Vault Agent injector example |
| **M438** | External Secrets Operator recipe |
| **M439** | SOPS <-> .katame interop bridge |
| **M440** | age <-> .katame conversion tool |
| **M441** | gpg <-> .katame conversion tool |
| **M442** | zipcrypto migration detector |
| **M443** | Malware scan hook after decrypt (optional) |
| **M444** | DLP label propagation into manifest |
| **M445** | Microsoft Purview labeling experiment |
| **M446** | Google Drive label experiment |
| **M447** | SharePoint retention hold recipe |
| **M448** | Legal hold audit export |
| **M449** | eDiscovery packager for sealed dirs |
| **M450** | Chain of custody report template |
| **M451** | Court-admissible hash export format |
| **M452** | RFC 3161 timestamp authority support |
| **M453** | OpenTimestamps alternate anchor |
| **M454** | Bitcoin OP_RETURN anchor option (opt-in) |
| **M455** | Certificate Transparency style monitor for Rekor |
| **M456** | Gossip protocol for repo root distribution |
| **M457** | Peer-to-peer bundle sync |
| **M458** | IPFS CID field on manifest |
| **M459** | Torrent magnet for large batch outputs |
| **M460** | CDN cache headers for public verify bundles |
| **M461** | Signed URL helper for cloud outs |
| **M462** | Presigned POST policy for protect uploads |
| **M463** | Bandwidth throttle for batch cloud |
| **M464** | Resume interrupted batch-protect |
| **M465** | Checkpoint journal for long batch jobs |
| **M466** | Distributed batch workers over NATS |
| **M467** | Redis lock for concurrent rotate-key |
| **M468** | etcd lease for serve leader election |
| **M469** | Active-active serve with sticky KMS cache |
| **M470** | Read replicas for repo verify |
| **M471** | Sharding strategy for huge repos |
| **M472** | Bloom filter for repo membership |
| **M473** | Compressed entries.jsonl option |
| **M474** | SQLite backend alternative to JSONL repo |
| **M475** | Postgres backend for enterprise repo |
| **M476** | S3 as repo object store |
| **M477** | Git-backed repo storage experiment |
| **M478** | FUSE filesystem mounting .katame dirs |
| **M479** | 9p serve of decrypted views (dangerous; gated) |
| **M480** | Windows ProjFS provider experiment |
| **M481** | macOS File Provider extension experiment |
| **M482** | rclone backend for .katame |
| **M483** | restic-compatible repository bridge |
| **M484** | borgbackup interop notes |
| **M485** | Tarsnap-like key hierarchy design |
| **M486** | Hierarchical DEKs for directory trees |
| **M487** | Per-file AAD binding to relative path |
| **M488** | Path traversal safe extract helper |
| **M489** | Unicode normalization for paths |
| **M490** | Case-folding policy on Windows/macOS |
| **M491** | Symlink handling policy in batch |
| **M492** | Hardlink dedup in batch-protect |
| **M493** | Sparse file support |
| **M494** | Extended attributes preservation |
| **M495** | ACL preservation cross-platform |
| **M496** | Owner/group mapping tables |
| **M497** | SELinux context notes |
| **M498** | AppArmor profile for katame serve |
| **M499** | seccomp.json for Docker serve |
| **M500** | Kubernetes NetworkPolicy examples |
| **M501** | Istio AuthorizationPolicy examples |
| **M502** | OPA Gatekeeper constraint for images |
| **M503** | Kyverno policy for cosign verify |
| **M504** | Supply chain levels for SLSA 3+ |
| **M505** | in-toto layout for release pipeline |
| **M506** | Witness run for build attestation |
| **M507** | GUAC ingestion of katame SBOMs |
| **M508** | Dependency-Track upload job |
| **M509** | OSV scanner in CI |
| **M510** | Nancy / govulncheck fail-on-high |
| **M511** | Trivy fs scan in CI |
| **M512** | Grype scan of release artifacts |
| **M513** | Scorecard checks enabling |
| **M514** | Allstar policy for org |
| **M515** | CodeQL Go queries custom pack |
| **M516** | Semgrep rules for crypto misuse |
| **M517** | gosec in lint pipeline |
| **M518** | staticcheck extras |
| **M519** | fieldalignment fixes |
| **M520** | shadow variable lint |
| **M521** | errname lint |
| **M522** | wrapcheck lint |
| **M523** | nilnil lint |
| **M524** | exhaustive switch lint for algos |
| **M525** | revive style guide |
| **M526** | gofumpt in CI |
| **M527** | goimports local prefix |
| **M528** | editorconfig enforcement |
| **M529** | pre-commit hooks repo config |
| **M530** | lefthook alternative |
| **M531** | Commit message convention (conventional commits) |
| **M532** | Release please / changesets |
| **M533** | Changelog automation from milestones |
| **M534** | Public roadmap board sync from milestones_backlog |
| **M535** | Milestone burndown metrics |
| **M536** | Contributor onboarding tutorial |
| **M537** | Good-first-issue labeling for M31+ |
| **M538** | Security.txt and vulnerability disclosure policy |
| **M539** | Bug bounty scope document |
| **M540** | Safe harbor research guidelines |
| **M541** | Code of conduct enforcement ladder |
| **M542** | DCO vs CLA decision |
| **M543** | Reproducible Windows builds with zig cc |
| **M544** | Cross sysroot for c-shared |
| **M545** | cgo-free PKCS#11 via process RPC |
| **M546** | cgo-free PIV via ykman CLI adapter |
| **M547** | SoftHSM2 docker compose for CI |
| **M548** | YubiKey CI with hardware farm (opt-in) |
| **M549** | Fake HSM for unit tests |
| **M550** | Deterministic mock clock for policy time windows |
| **M551** | Injectible entropy for tests |
| **M552** | Record/replay HTTP for cloud SDKs |
| **M553** | VCR cassettes for Rekor |
| **M554** | Pact contracts for serve API |
| **M555** | Postman collection export from OpenAPI |
| **M556** | Insomnia collection |
| **M557** | HTTPie cookbook |
| **M558** | curl recipes in docs/recipes.md expansion |
| **M559** | PowerShell module wrapper |
| **M560** | Bash completion install script |
| **M561** | fish vendor_completions.d packaging |
| **M562** | Nushell completions |
| **M563** | Elvish completions |
| **M564** | Fig/Amazon Q autocomplete spec |
| **M565** | VS Code extension for .katame inspect |
| **M566** | JetBrains plugin for .katame |
| **M567** | Vim/Neovim plugin |
| **M568** | Emacs package |
| **M569** | Language server for policy YAML |
| **M570** | Schema store entry for policy.yaml |
| **M571** | JSON Schema for .katame.yaml config |
| **M572** | cue / rego dual policy authoring |
| **M573** | Cedar policy backend experiment |
| **M574** | Polar (Oso) backend experiment |
| **M575** | ABAC attributes from cloud IAM |
| **M576** | RBAC role templates |
| **M577** | Break-glass audit alerts to PagerDuty |
| **M578** | Slack notifier for policy deny |
| **M579** | Teams webhook for seal mismatch |
| **M580** | Email digests from report |
| **M581** | CSV to SIEM mapping guide |
| **M582** | Splunk TA for audit JSONL |
| **M583** | Elastic ingest pipeline |
| **M584** | Loki labels for serve logs |
| **M585** | Datadog custom metrics mapping |
| **M586** | New Relic integration |
| **M587** | OpenSearch dashboard |
| **M588** | Chronicle / SecOps parser |
| **M589** | STIX/TAXII export of incidents |
| **M590** | MISP event from seal failure |
| **M591** | Sigma rules for katame audit |
| **M592** | YARA rules for .katame magic |
| **M593** | File magic / MIME database entry |
| **M594** | Windows Defender ASR notes |
| **M595** | CrowdStrike custom IOA examples |
| **M596** | osquery table for .katame files |
| **M597** | Fleet policy for encrypted exports |
| **M598** | Jamf profile for macOS keychain items |
| **M599** | Intune Win32 app packaging |
| **M600** | SCCM detection method |
| **M601** | Ansible role |
| **M602** | Chef cookbook |
| **M603** | Puppet module |
| **M604** | Salt state |
| **M605** | Terraform provider for serve |
| **M606** | Pulumi component |
| **M607** | Crossplane XRD |
| **M608** | Operator SDK CRD for KatameServe |
| **M609** | Helm chart with values schema |
| **M610** | Kustomize base overlays |
| **M611** | Skaffold profile |
| **M612** | Tiltfile for local serve |
| **M613** | Devcontainer feature |
| **M614** | GitHub Codespaces setup |
| **M615** | Gitpod config |
| **M616** | Docker bake hcl |
| **M617** | Buildx cache from GHCR |
| **M618** | Smoke test image on every PR |
| **M619** | Canary deploy workflow |
| **M620** | Blue/green serve rollout guide |
| **M621** | Chaos Mesh experiments |
| **M622** | Litmus chaos for KMS outage |
| **M623** | Game day runbook |
| **M624** | On-call cheatsheet |
| **M625** | RTO/RPO targets for repo |
| **M626** | Backup verification of sealed dirs |
| **M627** | DR drill checklist |
| **M628** | Multi-region active/passive serve |
| **M629** | Cold standby binary update procedure |
| **M630** | Key ceremony runbook for root signing |
| **M631** | Quarterly key rotation calendar |
| **M632** | Compromised-key response playbook |
| **M633** | CVE response SLA |
| **M634** | Customer advisory template |
| **M635** | Status page components list |
| **M636** | Public metrics of release health |
| **M637** | Community Discord/Slack bridge |
| **M638** | Office hours calendar |
| **M639** | Annual crypto design review |
| **M640** | External audit engagement checklist |
| **M641** | Penetration test scope for serve |
| **M642** | Red team exercise notes |
| **M643** | Purple team audit log scenarios |
| **M644** | Tabletop exercise: lost Shamir shares |
| **M645** | Tabletop exercise: KMS region outage |
| **M646** | Tabletop exercise: Rekor downtime |
| **M647** | Tabletop exercise: malicious plugin |
| **M648** | Insurance questionnaire answers |
| **M649** | SOC2 control mapping |
| **M650** | ISO 27001 Annex A mapping |
| **M651** | NIST 800-53 control mapping |
| **M652** | CIS benchmarks mapping |
| **M653** | PCI DSS key management notes |
| **M654** | HIPAA guidance for PHI in .katame |
| **M655** | GDPR encryption of personal data notes |
| **M656** | FedRAMP boundary discussion |
| **M657** | ITAR/export compliance for PQ crypto |
| **M658** | EAR classification notes |
| **M659** | Documentation translations (es) |
| **M660** | Documentation translations (pt) |
| **M661** | Documentation translations (de) |
| **M662** | Documentation translations (fr) |
| **M663** | Documentation translations (ja) |
| **M664** | Documentation translations (zh) |
| **M665** | RTL-friendly docs layout |
| **M666** | Accessibility review of menu TUI |
| **M667** | Screen reader labels for huh forms |
| **M668** | High-contrast terminal theme |
| **M669** | Colorblind-safe status glyphs |
| **M670** | Non-English error message catalogs |
| **M671** | ICU message format for errors |
| **M672** | Pluralization in report summaries |
| **M673** | Locale-aware timestamp display |
| **M674** | Timezone display preference in config |
| **M675** | Calendar integration for seal schedules |
| **M676** | Cron helper for seal + repo add |
| **M677** | GitHub Action composite for protect+sign |
| **M678** | GitLab CI component |
| **M679** | CircleCI orb |
| **M680** | Azure Pipelines task |
| **M681** | Jenkins shared library |
| **M682** | Tekton task catalog entry |
| **M683** | Argo Workflows template |
| **M684** | Airflow operator |
| **M685** | Prefect task |
| **M686** | Dagster asset |
| **M687** | dbt post-hook encryption example |
| **M688** | Spark UDF feasibility study |
| **M689** | Flink operator notes |
| **M690** | Kafka Connect SMT for .katame |
| **M691** | Pulsar function example |
| **M692** | NATS JetStream consumer example |
| **M693** | SQS worker pattern |
| **M694** | SNS fanout of attest events |
| **M695** | EventBridge rule examples |
| **M696** | CloudEvents binding for audit entries |
| **M697** | AsyncAPI for serve streaming future |
| **M698** | GraphQL facade experiment |
| **M699** | JSON:API envelope alternative |
| **M700** | Protobuf schema for manifest |
| **M701** | FlatBuffers experiment |
| **M702** | CapnProto experiment |
| **M703** | MessagePack manifest option |
| **M704** | CBOR manifest option |
| **M705** | Canonical CBOR for signatures |
| **M706** | JWS wrapping of signature.sig |
| **M707** | COSE_Sign1 alternative signature format |
| **M708** | CMS/PKCS7 detached signature option |
| **M709** | XMLDSig for enterprise interop |
| **M710** | PDF signing recipe using .katame sidecar |
| **M711** | Office document protect recipe |
| **M712** | DICOM medical image protect recipe |
| **M713** | GeoTIFF protect recipe |
| **M714** | Parquet columnar protect strategy |
| **M715** | ORC/Avro notes |
| **M716** | Arrow flight encryption notes |
| **M717** | ML model weight protect recipe |
| **M718** | ONNX bundle protect |
| **M719** | Docker layer protect experiment |
| **M720** | OCI artifact attach of .katame |
| **M721** | Sigstore cosign attach cookbook |
| **M722** | Notation (notaryv2) interop |
| **M723** | in-toto link for protect step |
| **M724** | Witness keyless for CLI in CI |
| **M725** | Git commit signing with katame keys |
| **M726** | GPG agent forwarding notes |
| **M727** | SSH certificate CA using katame keys |
| **M728** | SPIRE upstream authority experiment |
| **M729** | Cert-manager CSR signer plugin |
| **M730** | Smallstep CA integration |
| **M731** | Hashicorp Vault PKI recipe |
| **M732** | AWS ACM PCA recipe |
| **M733** | Google CAS recipe |
| **M734** | Azure Key Vault certs recipe |
| **M735** | Let's Encrypt for serve TLS automation |
| **M736** | ACME HTTP-01 sidecar |
| **M737** | mTLS mesh with linkerd |
| **M738** | mTLS mesh with consul connect |
| **M739** | Zero-trust access (Cloudflare/Teleport) recipe |
| **M740** | Tailscale serve funnel example |
| **M741** | WireGuard sidecar pattern |
| **M742** | SSH reverse tunnel for airgap serve |
| **M743** | YubiHSM2 support beyond PIV |
| **M744** | CloudHSM AWS provider |
| **M745** | CloudHSM GCP provider |
| **M746** | Azure Dedicated HSM provider |
| **M747** | Thales Luna PKCS#11 recipe |
| **M748** | Utimaco PKCS#11 recipe |
| **M749** | nCipher/Entrust PKCS#11 recipe |
| **M750** | AWS Payment Cryptography provider |
| **M751** | IBM Hyper Protect notes |
| **M752** | SEV-SNP sealed DEK experiment |
| **M753** | TDX sealed DEK experiment |
| **M754** | SGX enclave protect experiment |
| **M755** | TrustZone notes |
| **M756** | Android StrongBox binding |
| **M757** | iOS Secure Enclave binding |
| **M758** | Flutter plugin feasibility |
| **M759** | React Native turbo module |
| **M760** | Kotlin Multiplatform FFI |
| **M761** | Swift Package wrapping C ABI |
| **M762** | Rust crate via bindgen |
| **M763** | Zig package via @cImport |
| **M764** | Nim wrapper |
| **M765** | Crystal binding |
| **M766** | Haskell FFI |
| **M767** | OCaml ctypes |
| **M768** | Erlang NIF |
| **M769** | Elixir rustler alternative using C ABI |
| **M770** | PHP FFI extension |
| **M771** | Perl XS thin wrapper |
| **M772** | R package Rcpp for analysis of manifests |
| **M773** | Julia ccall wrapper |
| **M774** | MATLAB loadlibrary example |
| **M775** | LabVIEW DLL call notes |
| **M776** | Excel add-in for inspect (joke-tier but tracked) |
| **M777** | Power BI connector for audit reports |
| **M778** | Tableau Web Data Connector |
| **M779** | Observable Plot dashboard of report JSON |
| **M780** | Streamlit app for local decrypt UI |
| **M781** | Gradio demo |
| **M782** | Hugging Face Space demo (public keys only) |
| **M783** | Replit template |
| **M784** | Glitch template |
| **M785** | CodeSandbox WASM demo |
| **M786** | StackBlitz WASM demo |
| **M787** | Benchmark published on website |
| **M788** | Public corpus of encrypted samples |
| **M789** | CTF challenges using .katame |
| **M790** | University course lab materials |
| **M791** | Workshop slide deck |
| **M792** | Conference talk outline |
| **M793** | Blog series: Streaming AEAD deep dive |
| **M794** | Blog series: PQ migration path |
| **M795** | Blog series: Policy fail-closed lessons |
| **M796** | Podcast interview FAQ |
| **M797** | Changelog video automation |
| **M798** | Twitter/X release bot |
| **M799** | Mastodon release bot |
| **M800** | Bluesky release bot |
| **M801** | RSS feed for releases |
| **M802** | Atom feed for security advisories |
| **M803** | Security advisory PGP key |
| **M804** | Keyserver upload automation |
| **M805** | WKD for release identity |
| **M806** | DNS CAA for katame.example |
| **M807** | Brand assets / logo SVG |
| **M808** | Sticker sheet for conferences |
| **M809** | Swag store (optional) |
| **M810** | Contributor ladder |
| **M811** | Maintainer rota |
| **M812** | Triage bot for issues |
| **M813** | Stale bot config |
| **M814** | Issue forms for bug/feature |
| **M815** | PR template with test plan |
| **M816** | CODEOWNERS for crypto paths |
| **M817** | Branch protection documentation |
| **M818** | Required review for internal/crypto |
| **M819** | Signed commits required |
| **M820** | Merge queue enablement |
| **M821** | Auto-assign reviewers |
| **M822** | Danger JS for PR checks |
| **M823** | Semantic PR title lint |
| **M824** | Size labeler for PRs |
| **M825** | Release drafter config |
| **M826** | Milestone auto-close on tag |
| **M827** | Project board automation |
| **M828** | OKR tracking for M31-M100 |
| **M829** | Quarterly prune of wontfix milestones |
| **M830** | Archive completed milestones annually |
| **M831** | Add parallel encrypt correctness test for cipher aes-256-gcm |
| **M832** | Document performance profile for cipher aes-256-gcm |
| **M833** | Add parallel encrypt correctness test for cipher chacha20-poly1305 |
| **M834** | Document performance profile for cipher chacha20-poly1305 |
| **M835** | Add parallel encrypt correctness test for cipher xchacha20-poly1305 |
| **M836** | Document performance profile for cipher xchacha20-poly1305 |
| **M837** | Add parallel encrypt correctness test for cipher aes-256-gcm-siv |
| **M838** | Document performance profile for cipher aes-256-gcm-siv |
| **M839** | Add protect/decrypt round-trip with hash-algo sha256 |
| **M840** | Add fuzz corpus seed for hash sha256 |
| **M841** | Add protect/decrypt round-trip with hash-algo sha512 |
| **M842** | Add fuzz corpus seed for hash sha512 |
| **M843** | Add protect/decrypt round-trip with hash-algo sha3-256 |
| **M844** | Add fuzz corpus seed for hash sha3-256 |
| **M845** | Add protect/decrypt round-trip with hash-algo sha3-512 |
| **M846** | Add fuzz corpus seed for hash sha3-512 |
| **M847** | Add protect/decrypt round-trip with hash-algo blake2b-256 |
| **M848** | Add fuzz corpus seed for hash blake2b-256 |
| **M849** | Add protect/decrypt round-trip with hash-algo blake2b-512 |
| **M850** | Add fuzz corpus seed for hash blake2b-512 |
| **M851** | Add protect/decrypt round-trip with hash-algo blake3 |
| **M852** | Add fuzz corpus seed for hash blake3 |
| **M853** | Bench KDF argon2id on CI reference hardware |
| **M854** | Document recommended params for KDF argon2id |
| **M855** | Bench KDF scrypt on CI reference hardware |
| **M856** | Document recommended params for KDF scrypt |
| **M857** | Bench KDF pbkdf2-sha256 on CI reference hardware |
| **M858** | Document recommended params for KDF pbkdf2-sha256 |
| **M859** | Add golden signature vector for ed25519 |
| **M860** | Add verify negative test for ed25519 with wrong key |
| **M861** | Add golden signature vector for ed448 |
| **M862** | Add verify negative test for ed448 with wrong key |
| **M863** | Add golden signature vector for ecdsa-p256 |
| **M864** | Add verify negative test for ecdsa-p256 with wrong key |
| **M865** | Add golden signature vector for ecdsa-p384 |
| **M866** | Add verify negative test for ecdsa-p384 with wrong key |
| **M867** | Add golden signature vector for rsa-pss-2048 |
| **M868** | Add verify negative test for rsa-pss-2048 with wrong key |
| **M869** | Add golden signature vector for rsa-pss-4096 |
| **M870** | Add verify negative test for rsa-pss-4096 with wrong key |
| **M871** | Add golden signature vector for ml-dsa-65 |
| **M872** | Add verify negative test for ml-dsa-65 with wrong key |
| **M873** | Add golden signature vector for ml-dsa-87 |
| **M874** | Add verify negative test for ml-dsa-87 with wrong key |
| **M875** | Add golden signature vector for slh-dsa-128s |
| **M876** | Add verify negative test for slh-dsa-128s with wrong key |
| **M877** | Add multi-recipient test including scheme x25519-aes-256-gcm |
| **M878** | Document PEM block types for x25519-aes-256-gcm |
| **M879** | Add multi-recipient test including scheme ecies-p256 |
| **M880** | Document PEM block types for ecies-p256 |
| **M881** | Add multi-recipient test including scheme rsa-oaep-2048 |
| **M882** | Document PEM block types for rsa-oaep-2048 |
| **M883** | Add multi-recipient test including scheme rsa-oaep-4096 |
| **M884** | Document PEM block types for rsa-oaep-4096 |
| **M885** | Add multi-recipient test including scheme ml-kem-768 |
| **M886** | Document PEM block types for ml-kem-768 |
| **M887** | Add multi-recipient test including scheme ml-kem-1024 |
| **M888** | Document PEM block types for ml-kem-1024 |
| **M889** | Add retry classification matrix for cloud s3 |
| **M890** | Add doctor check for s3 credentials |
| **M891** | Add retry classification matrix for cloud gcs |
| **M892** | Add doctor check for gcs credentials |
| **M893** | Add retry classification matrix for cloud azure |
| **M894** | Add doctor check for azure credentials |
| **M895** | Add retry classification matrix for cloud https |
| **M896** | Add doctor check for https credentials |
| **M897** | Expose rewrap in Go SDK |
| **M898** | Expose rotate-key in Go SDK |
| **M899** | Document error mapping for Go SDK |
| **M900** | Add smoke test for compress=zstd in Go SDK |
| **M901** | Add smoke test for hybrid recipient in Go SDK |
| **M902** | Expose rewrap in C ABI |
| **M903** | Expose rotate-key in C ABI |
| **M904** | Document error mapping for C ABI |
| **M905** | Add smoke test for compress=zstd in C ABI |
| **M906** | Add smoke test for hybrid recipient in C ABI |
| **M907** | Expose rewrap in Python |
| **M908** | Expose rotate-key in Python |
| **M909** | Document error mapping for Python |
| **M910** | Add smoke test for compress=zstd in Python |
| **M911** | Add smoke test for hybrid recipient in Python |
| **M912** | Expose rewrap in TypeScript/WASM |
| **M913** | Expose rotate-key in TypeScript/WASM |
| **M914** | Document error mapping for TypeScript/WASM |
| **M915** | Add smoke test for compress=zstd in TypeScript/WASM |
| **M916** | Add smoke test for hybrid recipient in TypeScript/WASM |
| **M917** | Expose rewrap in serve HTTP |
| **M918** | Expose rotate-key in serve HTTP |
| **M919** | Document error mapping for serve HTTP |
| **M920** | Add smoke test for compress=zstd in serve HTTP |
| **M921** | Add smoke test for hybrid recipient in serve HTTP |
| **M922** | Add structured logging field op=protect |
| **M923** | Add audit coverage assertion for protect |
| **M924** | Add --help example block for protect |
| **M925** | Add JSON schema for protect --json output |
| **M926** | Add structured logging field op=decrypt |
| **M927** | Add audit coverage assertion for decrypt |
| **M928** | Add --help example block for decrypt |
| **M929** | Add JSON schema for decrypt --json output |
| **M930** | Add structured logging field op=inspect |
| **M931** | Add audit coverage assertion for inspect |
| **M932** | Add --help example block for inspect |
| **M933** | Add JSON schema for inspect --json output |
| **M934** | Add structured logging field op=diff |
| **M935** | Add audit coverage assertion for diff |
| **M936** | Add --help example block for diff |
| **M937** | Add JSON schema for diff --json output |
| **M938** | Add structured logging field op=sign |
| **M939** | Add audit coverage assertion for sign |
| **M940** | Add --help example block for sign |
| **M941** | Add JSON schema for sign --json output |
| **M942** | Add structured logging field op=verify |
| **M943** | Add audit coverage assertion for verify |
| **M944** | Add --help example block for verify |
| **M945** | Add JSON schema for verify --json output |
| **M946** | Add structured logging field op=verify-integrity |
| **M947** | Add audit coverage assertion for verify-integrity |
| **M948** | Add --help example block for verify-integrity |
| **M949** | Add JSON schema for verify-integrity --json output |
| **M950** | Add structured logging field op=rewrap |
| **M951** | Add audit coverage assertion for rewrap |
| **M952** | Add --help example block for rewrap |
| **M953** | Add JSON schema for rewrap --json output |
| **M954** | Add structured logging field op=rotate-key |
| **M955** | Add audit coverage assertion for rotate-key |
| **M956** | Add --help example block for rotate-key |
| **M957** | Add JSON schema for rotate-key --json output |
| **M958** | Add structured logging field op=add-recipient |
| **M959** | Add audit coverage assertion for add-recipient |
| **M960** | Add --help example block for add-recipient |
| **M961** | Add JSON schema for add-recipient --json output |
| **M962** | Add structured logging field op=remove-recipient |
| **M963** | Add audit coverage assertion for remove-recipient |
| **M964** | Add --help example block for remove-recipient |
| **M965** | Add JSON schema for remove-recipient --json output |
| **M966** | Add structured logging field op=attest |
| **M967** | Add audit coverage assertion for attest |
| **M968** | Add --help example block for attest |
| **M969** | Add JSON schema for attest --json output |
| **M970** | Add structured logging field op=seal |
| **M971** | Add audit coverage assertion for seal |
| **M972** | Add --help example block for seal |
| **M973** | Add JSON schema for seal --json output |
| **M974** | Add structured logging field op=verify-seal |
| **M975** | Add audit coverage assertion for verify-seal |
| **M976** | Add --help example block for verify-seal |
| **M977** | Add JSON schema for verify-seal --json output |
| **M978** | Add structured logging field op=repo |
| **M979** | Add audit coverage assertion for repo |
| **M980** | Add --help example block for repo |
| **M981** | Add JSON schema for repo --json output |
| **M982** | Add structured logging field op=serve |
| **M983** | Add audit coverage assertion for serve |
| **M984** | Add --help example block for serve |
| **M985** | Add JSON schema for serve --json output |
| **M986** | Add structured logging field op=policy |
| **M987** | Add audit coverage assertion for policy |
| **M988** | Add --help example block for policy |
| **M989** | Add JSON schema for policy --json output |
| **M990** | Add structured logging field op=doctor |
| **M991** | Add audit coverage assertion for doctor |
| **M992** | Add --help example block for doctor |
| **M993** | Add JSON schema for doctor --json output |
| **M994** | Add structured logging field op=report |
| **M995** | Add audit coverage assertion for report |
| **M996** | Add --help example block for report |
| **M997** | Add JSON schema for report --json output |
| **M998** | Add structured logging field op=batch-protect |
| **M999** | Add audit coverage assertion for batch-protect |
| **M1000** | Add --help example block for batch-protect |
| **M1001** | Add JSON schema for batch-protect --json output |
| **M1002** | Add structured logging field op=batch-decrypt |
| **M1003** | Add audit coverage assertion for batch-decrypt |
| **M1004** | Add --help example block for batch-decrypt |
| **M1005** | Add JSON schema for batch-decrypt --json output |
| **M1006** | Add structured logging field op=batch-inspect |
| **M1007** | Add audit coverage assertion for batch-inspect |
| **M1008** | Add --help example block for batch-inspect |
| **M1009** | Add JSON schema for batch-inspect --json output |
| **M1010** | Add structured logging field op=split-key |
| **M1011** | Add audit coverage assertion for split-key |
| **M1012** | Add --help example block for split-key |
| **M1013** | Add JSON schema for split-key --json output |
| **M1014** | Add structured logging field op=combine-key |
| **M1015** | Add audit coverage assertion for combine-key |
| **M1016** | Add --help example block for combine-key |
| **M1017** | Add JSON schema for combine-key --json output |
| **M1018** | Add structured logging field op=hash |
| **M1019** | Add audit coverage assertion for hash |
| **M1020** | Add --help example block for hash |
| **M1021** | Add JSON schema for hash --json output |
| **M1022** | Add structured logging field op=keygen |
| **M1023** | Add audit coverage assertion for keygen |
| **M1024** | Add --help example block for keygen |
| **M1025** | Add JSON schema for keygen --json output |
| **M1026** | Add structured logging field op=config |
| **M1027** | Add audit coverage assertion for config |
| **M1028** | Add --help example block for config |
| **M1029** | Add JSON schema for config --json output |
| **M1030** | Add structured logging field op=menu |
| **M1031** | Add audit coverage assertion for menu |
| **M1032** | Add --help example block for menu |
| **M1033** | Add JSON schema for menu --json output |
| **M1034** | Add structured logging field op=completion |
| **M1035** | Add audit coverage assertion for completion |
| **M1036** | Add --help example block for completion |
| **M1037** | Add JSON schema for completion --json output |
| **M1038** | Add structured logging field op=audit |
| **M1039** | Add audit coverage assertion for audit |
| **M1040** | Add --help example block for audit |
| **M1041** | Add JSON schema for audit --json output |

Total backlog entries in this file: 1011.
