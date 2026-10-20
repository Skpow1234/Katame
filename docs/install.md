# Install

Pre-built binaries are published for **Linux** (amd64, arm64), **macOS** (amd64, arm64), and **Windows** (amd64, arm64) on [Releases](https://github.com/Skpow1234/Vaultpack/releases). Each release includes checksums and SBOMs.

Requires **Go 1.26+** when building from source (see `go.mod` / `toolchain`).

## One-command install (curl)

Replace `VERSION` with a tag (e.g. `v1.0.0`) or use `latest`:

```bash
# Linux (amd64)
curl -sL https://github.com/Skpow1234/Vaultpack/releases/latest/download/vaultpack-linux-amd64.tar.gz | tar xz
sudo mv vaultpack /usr/local/bin/

# Linux (arm64)
curl -sL https://github.com/Skpow1234/Vaultpack/releases/latest/download/vaultpack-linux-arm64.tar.gz | tar xz
sudo mv vaultpack /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/Skpow1234/Vaultpack/releases/latest/download/vaultpack-darwin-arm64.tar.gz | tar xz
sudo mv vaultpack /usr/local/bin/

# macOS (Intel)
curl -sL https://github.com/Skpow1234/Vaultpack/releases/latest/download/vaultpack-darwin-amd64.tar.gz | tar xz
sudo mv vaultpack /usr/local/bin/
```

**Windows:** download `vaultpack-windows-amd64.zip` or `vaultpack-windows-arm64.zip` from Releases, extract, and add the folder to your `PATH`.

## Package managers

```bash
# Homebrew (macOS / Linux) — after tap is set up
brew install Skpow1234/tap/vaultpack

# Scoop (Windows)
scoop bucket add vaultpack https://github.com/Skpow1234/scoop-bucket
scoop install vaultpack

# Winget (Windows) — after PR is merged into microsoft/winget-pkgs
winget install Skpow1234.Vaultpack

# Chocolatey (Windows)
choco install vaultpack

# Debian / Ubuntu (.deb)
curl -sLO https://github.com/Skpow1234/Vaultpack/releases/latest/download/vaultpack_<version>_linux_amd64.deb
sudo dpkg -i vaultpack_<version>_linux_amd64.deb

# Red Hat / Fedora (.rpm)
curl -sLO https://github.com/Skpow1234/Vaultpack/releases/latest/download/vaultpack_<version>_linux_amd64.rpm
sudo rpm -ivh vaultpack_<version>_linux_amd64.rpm

# Alpine (.apk)
curl -sLO https://github.com/Skpow1234/Vaultpack/releases/latest/download/vaultpack_<version>_linux_amd64.apk
sudo apk add --allow-untrusted vaultpack_<version>_linux_amd64.apk
```

`.deb` / `.rpm` / `.apk` artifacts are produced by [nfpm](https://nfpm.goreleaser.com/) in the GoReleaser pipeline. A hosted apt/yum repo (Cloudsmith) can be enabled with `CLOUDSMITH_API_KEY` in release secrets.

### Maintainer setup (Scoop / Winget / Chocolatey)

| Channel | One-time setup | GitHub secret | User install |
| ------- | -------------- | ------------- | ------------ |
| **Scoop** | Create `Skpow1234/scoop-bucket` | `SCOOP_BUCKET_GITHUB_TOKEN` | `scoop bucket add vaultpack …` then `scoop install vaultpack` |
| **Winget** | Fork [microsoft/winget-pkgs](https://github.com/microsoft/winget-pkgs) | `WINGET_PUBLISH_TOKEN` | `winget install Skpow1234.Vaultpack` |
| **Chocolatey** | Register at chocolatey.org | `CHOCOLATEY_API_KEY` | `choco install vaultpack` |

## Verify checksums

```bash
curl -sLO https://github.com/Skpow1234/Vaultpack/releases/latest/download/checksums-sha256.txt
curl -sLO https://github.com/Skpow1234/Vaultpack/releases/latest/download/vaultpack-linux-amd64.tar.gz
sha256sum -c checksums-sha256.txt --ignore-missing
```

On macOS use `shasum -a 256 -c checksums-sha256.txt --ignore-missing`.

## Verify cosign signatures (Sigstore keyless)

```bash
ARCHIVE=vaultpack-linux-amd64.tar.gz
BASE=https://github.com/Skpow1234/Vaultpack/releases/latest/download
curl -sLO $BASE/$ARCHIVE
curl -sLO $BASE/$ARCHIVE.sig
curl -sLO $BASE/$ARCHIVE.pem

cosign verify-blob \
  --certificate $ARCHIVE.pem \
  --signature $ARCHIVE.sig \
  --certificate-identity-regexp 'https://github.com/Skpow1234/Vaultpack' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  $ARCHIVE
```

## With `go install`

```bash
go install github.com/Skpow1234/Vaultpack/cmd/vaultpack@latest
```

## From source

```bash
git clone https://github.com/Skpow1234/Vaultpack.git
cd Vaultpack
go build -o bin/vaultpack ./cmd/vaultpack
```

## Docker

Images target **linux/amd64** and **linux/arm64** (distroless runtime):

```bash
docker run --rm -v "$PWD:/work" ghcr.io/Skpow1234/vaultpack:latest protect --in /work/config.json

docker build -t vaultpack .
docker run --rm -v "$PWD:/work" vaultpack protect --in /work/config.json
```

Multi-arch: `make docker-buildx` or `docker buildx build --platform linux/amd64,linux/arm64 -t vaultpack:tag .`

See also [Deployment](deployment.md).
