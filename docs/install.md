# Install

Pre-built binaries are published for **Linux** (amd64, arm64), **macOS** (amd64, arm64), and **Windows** (amd64, arm64) on [Releases](https://github.com/Skpow1234/Katame/releases). Each release includes checksums and SBOMs.

Requires **Go 1.26+** when building from source (see `go.mod` / `toolchain`).

## One-command install (curl)

Replace `VERSION` with a tag (e.g. `v1.0.0`) or use `latest`:

```bash
# Linux (amd64)
curl -sL https://github.com/Skpow1234/Katame/releases/latest/download/katame-linux-amd64.tar.gz | tar xz
sudo mv katame /usr/local/bin/

# Linux (arm64)
curl -sL https://github.com/Skpow1234/Katame/releases/latest/download/katame-linux-arm64.tar.gz | tar xz
sudo mv katame /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/Skpow1234/Katame/releases/latest/download/katame-darwin-arm64.tar.gz | tar xz
sudo mv katame /usr/local/bin/

# macOS (Intel)
curl -sL https://github.com/Skpow1234/Katame/releases/latest/download/katame-darwin-amd64.tar.gz | tar xz
sudo mv katame /usr/local/bin/
```

**Windows:** download `katame-windows-amd64.zip` or `katame-windows-arm64.zip` from Releases, extract, and add the folder to your `PATH`.

## Package managers

```bash
# Homebrew (macOS / Linux) — after tap is set up
brew install Skpow1234/tap/katame

# Scoop (Windows)
scoop bucket add katame https://github.com/Skpow1234/scoop-bucket
scoop install katame

# Winget (Windows) — after PR is merged into microsoft/winget-pkgs
winget install Skpow1234.Katame

# Chocolatey (Windows)
choco install katame

# Debian / Ubuntu (.deb)
curl -sLO https://github.com/Skpow1234/Katame/releases/latest/download/katame_<version>_linux_amd64.deb
sudo dpkg -i katame_<version>_linux_amd64.deb

# Red Hat / Fedora (.rpm)
curl -sLO https://github.com/Skpow1234/Katame/releases/latest/download/katame_<version>_linux_amd64.rpm
sudo rpm -ivh katame_<version>_linux_amd64.rpm

# Alpine (.apk)
curl -sLO https://github.com/Skpow1234/Katame/releases/latest/download/katame_<version>_linux_amd64.apk
sudo apk add --allow-untrusted katame_<version>_linux_amd64.apk
```

`.deb` / `.rpm` / `.apk` artifacts are produced by [nfpm](https://nfpm.goreleaser.com/) in the GoReleaser pipeline. A hosted apt/yum repo (Cloudsmith) can be enabled with `CLOUDSMITH_API_KEY` in release secrets.

### Maintainer setup (Scoop / Winget / Chocolatey)

| Channel | One-time setup | GitHub secret | User install |
| ------- | -------------- | ------------- | ------------ |
| **Scoop** | Create `Skpow1234/scoop-bucket` | `SCOOP_BUCKET_GITHUB_TOKEN` | `scoop bucket add katame …` then `scoop install katame` |
| **Winget** | Fork [microsoft/winget-pkgs](https://github.com/microsoft/winget-pkgs) | `WINGET_PUBLISH_TOKEN` | `winget install Skpow1234.Katame` |
| **Chocolatey** | Register at chocolatey.org | `CHOCOLATEY_API_KEY` | `choco install katame` |

## Verify checksums

```bash
curl -sLO https://github.com/Skpow1234/Katame/releases/latest/download/checksums-sha256.txt
curl -sLO https://github.com/Skpow1234/Katame/releases/latest/download/katame-linux-amd64.tar.gz
sha256sum -c checksums-sha256.txt --ignore-missing
```

On macOS use `shasum -a 256 -c checksums-sha256.txt --ignore-missing`.

## Verify cosign signatures (Sigstore keyless)

```bash
ARCHIVE=katame-linux-amd64.tar.gz
BASE=https://github.com/Skpow1234/Katame/releases/latest/download
curl -sLO $BASE/$ARCHIVE
curl -sLO $BASE/$ARCHIVE.sig
curl -sLO $BASE/$ARCHIVE.pem

cosign verify-blob \
  --certificate $ARCHIVE.pem \
  --signature $ARCHIVE.sig \
  --certificate-identity-regexp 'https://github.com/Skpow1234/Katame' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  $ARCHIVE
```

## With `go install`

```bash
go install github.com/Skpow1234/Katame/cmd/katame@latest
```

## From source

```bash
git clone https://github.com/Skpow1234/Katame.git
cd Katame
go build -o bin/katame ./cmd/katame
```

## Docker

Images target **linux/amd64** and **linux/arm64** (distroless runtime):

```bash
docker run --rm -v "$PWD:/work" ghcr.io/Skpow1234/katame:latest protect --in /work/config.json

docker build -t katame .
docker run --rm -v "$PWD:/work" katame protect --in /work/config.json
```

Multi-arch: `make docker-buildx` or `docker buildx build --platform linux/amd64,linux/arm64 -t katame:tag .`

See also [Deployment](deployment.md).
