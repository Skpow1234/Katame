# Katame

[![CI](https://github.com/Skpow1234/Katame/actions/workflows/ci.yml/badge.svg)](https://github.com/Skpow1234/Katame/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A cross-platform CLI (and SDK) that encrypts, hashes, and signs data artifacts into portable `.katame` bundles.

One tool. One bundle. Encryption + integrity + authenticity.

## Quick start

```bash
# Encrypt (generates a key; default cipher AES-256-GCM)
katame protect --in config.json

# Decrypt
katame decrypt --in config.json.katame --out config.json --key config.json.key

# Hybrid (recipient public key)
katame keygen --out alice --algo x25519-aes-256-gcm
katame protect --in config.json --recipient alice.pub
katame decrypt --in config.json.katame --out config.json --privkey alice.key

# Inspect
katame inspect --in config.json.katame
```

## Install

```bash
# From GitHub Releases (Linux amd64 example)
curl -sL https://github.com/Skpow1234/Katame/releases/latest/download/katame-linux-amd64.tar.gz | tar xz
sudo mv katame /usr/local/bin/

# Or
go install github.com/Skpow1234/Katame/cmd/katame@latest
```

Full install options (packages, Docker, cosign verification): **[docs/install.md](docs/install.md)**  
How to deploy (binaries, containers, `serve`, sidecars): **[docs/deployment.md](docs/deployment.md)**

## Documentation

| Topic | Link |
| ----- | ---- |
| Docs index | [docs/README.md](docs/README.md) |
| CLI usage | [docs/usage.md](docs/usage.md) |
| HTTP API / endpoints | [docs/endpoints.md](docs/endpoints.md) |
| Service mode | [docs/serve.md](docs/serve.md) |
| Security model | [docs/security.md](docs/security.md) |
| Bundle format | [docs/design.md](docs/design.md) |
| SDK & bindings | [docs/sdk.md](docs/sdk.md) |
| Policy, audit, repo, cloud | [docs/](docs/README.md) |
| Recipes | [docs/recipes.md](docs/recipes.md) |
| Milestones | [docs/milestones.md](docs/milestones.md) · [backlog (1000+)](docs/milestones_backlog.md) |
| Development | [docs/development.md](docs/development.md) |

## Features at a glance

- AEAD encryption (AES-GCM, ChaCha20, XChaCha20, AES-GCM-SIV) with chunked streaming
- Password KDFs, hybrid multi-recipient, KMS wrap (AWS / GCP / Azure), Shamir shares
- Post-quantum: ML-KEM, ML-DSA, SLH-DSA
- Sign / verify, Sigstore Rekor transparency, SLSA-style attestations
- Policy (YAML/JSON/Rego), hash-chained audit, tamper-evident repo
- Cloud URIs (`az://`, `s3://`, `gs://`, HTTPS), batch ops, `katame serve` HTTP API
- Go SDK, C-shared library, Python and TypeScript/WASM bindings

## Bundle format

```text
artifact.katame
├── payload.bin        # ciphertext
├── manifest.json      # encryption params, hashes, metadata
├── signature.sig      # optional
└── provenance.json    # optional
```

Details: [docs/design.md](docs/design.md) · [docs/security.md](docs/security.md)

## License

[MIT](LICENSE)
