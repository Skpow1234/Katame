# Deployment

Ways to ship and run VaultPack, based on what the repo actually supports today.

## 1. Released binaries

Best for workstations and most servers.

- Download from [GitHub Releases](https://github.com/Skpow1234/Vaultpack/releases)
- Verify [checksums + cosign](install.md)
- Install via Homebrew / Scoop / Winget / Chocolatey / `.deb` / `.rpm` / `.apk` when channels are wired

CI: `.github/workflows/release.yml` + `.goreleaser.yaml` (multi-OS, SBOM, cosign keyless, optional Cloudsmith).

## 2. Docker / OCI

`Dockerfile` is a multi-stage build: `golang` builder → `gcr.io/distroless/static-debian12:nonroot`.

```bash
# Local
make docker-build
docker run --rm -v "$PWD:/work" vaultpack:latest protect --in /work/secret.txt

# Multi-arch
make docker-buildx
# or: docker buildx build --platform linux/amd64,linux/arm64 -t vaultpack:tag --push .

# Published image (when pushed)
docker run --rm -v "$PWD:/work" ghcr.io/Skpow1234/vaultpack:latest hash --in /work/file.bin
```

Entrypoint is `vaultpack`; pass subcommands as args. Mount a work volume — the image has no shell tooling beyond the binary.

## 3. HTTP service (`vaultpack serve`)

Run VaultPack as a long-lived API (see [serve.md](serve.md) and [endpoints.md](endpoints.md)).

### Sidecar (UNIX socket)

```bash
vaultpack serve --listen unix:/run/vaultpack.sock --auth-disabled \
  --policy /etc/vaultpack/policy.yaml \
  --audit-log /var/log/vaultpack-audit.jsonl
```

Peer containers/processes speak HTTP over the socket. Suitable for Kubernetes sidecars or systemd units sharing a private volume for the socket.

### Network listener (TLS)

```bash
vaultpack serve \
  --listen :8443 \
  --tls-cert /etc/ssl/vaultpack.crt \
  --tls-key /etc/ssl/vaultpack.key \
  --auth-token-file /run/secrets/vaultpack-token \
  --policy /etc/vaultpack/policy.yaml
```

Do **not** use `--auth-disabled` on a TCP listener in production. Prefer mTLS or OIDC for multi-tenant setups.

### Example systemd unit

```ini
[Unit]
Description=VaultPack API
After=network.target

[Service]
ExecStart=/usr/local/bin/vaultpack serve \
  --listen :8443 \
  --tls-cert /etc/vaultpack/tls.crt \
  --tls-key /etc/vaultpack/tls.key \
  --auth-token-file /etc/vaultpack/token \
  --policy /etc/vaultpack/policy.yaml \
  --audit-log /var/log/vaultpack-audit.jsonl
Restart=on-failure
User=vaultpack
Group=vaultpack

[Install]
WantedBy=multi-user.target
```

### Kubernetes sketch

There is no Helm chart in-repo. A minimal pattern:

1. Deploy the distroless image with `args: ["serve", "--listen", ":8443", ...]`
2. Mount TLS material and token from Secrets
3. Probe `GET /healthz`
4. Scrape `/metrics` only from the cluster network (or a metrics sidecar) — it is unauthenticated
5. Optional: second container sharing an `emptyDir` for `unix:/run/vaultpack.sock`

## 4. Library / process-local

- Link `pkg/vaultpack` in Go services
- Load `libvaultpack` from Python/TS via bindings
- Invoke the CLI from CI jobs (`protect` / `sign` / `verify` / `repo`)

## 5. Package-repo publishing (maintainers)

| Channel | Gate |
| ------- | ---- |
| GitHub Releases | Tag `v*.*.*` |
| Homebrew | Tap token in GoReleaser |
| Scoop / Winget | PATs in release secrets |
| Chocolatey | Windows job + API key |
| Cloudsmith apt/yum | `CLOUDSMITH_API_KEY` |

## Checklist

- [ ] Verify release signatures before promoting
- [ ] Set policy + audit log for serve and CI decrypt paths
- [ ] Restrict `/metrics` and never ship `--auth-disabled` on public ports
- [ ] Prefer KMS-wrapped DEKs over long-lived key files on disk
