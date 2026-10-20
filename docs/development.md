# Development

```bash
make build        # bin/vaultpack
make test         # race-enabled tests
make lint         # golangci-lint
make fmt
make vet
make vulncheck
make fuzz
make docker-build
make release-snapshot
```

Requires Go **1.26.x** (see `go.mod` `toolchain`). CI pins `GOTOOLCHAIN=local` with setup-go `1.26.10`.

## Project structure

```text
cmd/vaultpack/          # CLI entrypoint
cmd/vaultpack-c/        # C-shared library (cshared tag)
cmd/vaultpack-wasm/     # WASM build
pkg/vaultpack/          # Public Go SDK
internal/cli/           # Cobra commands
internal/crypto/        # AEAD, KDF, hybrid, PQ, SSS
internal/bundle/        # ZIP I/O, manifest
internal/vpackop/       # Shared protect/decrypt/rewrap ops
internal/kms/           # AWS / GCP / Azure / mock
internal/keysource/     # URI key loaders
internal/policy/        # YAML/JSON/Rego
internal/audit/         # Audit chain, Merkle seal, provenance
internal/serve/         # HTTP API
internal/repo/          # Tamper-evident repo
internal/transparency/  # Rekor / Fulcio
internal/plugin/        # Plugin discovery
internal/cloud/         # S3 / GCS / Azure / HTTP fetch
internal/config/        # Viper config
bindings/               # Python, TypeScript
examples/               # Plugin examples
docs/                   # This documentation
```

## Tests

~70+ `*_test.go` files including fuzz targets under `internal/crypto` and `internal/bundle`. CI runs unit tests on Linux/macOS/Windows, govulncheck, short fuzz, and cross-builds.
