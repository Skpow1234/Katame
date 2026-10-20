# Service mode (`vaultpack serve`)

Runs a long-running HTTP API exposing the SDK. The wire contract matches the C-shared library and WASM bindings: JSON in/out with an `{"ok": true|false, ...}` envelope.

Full route reference: [endpoints.md](endpoints.md) · OpenAPI: [openapi.yaml](openapi.yaml).

## Examples

```bash
# Local development only
vaultpack serve --listen :8443 --auth-disabled

# Bearer token + TLS
vaultpack serve \
  --listen :8443 \
  --auth-token "$(cat /etc/vaultpack/token)" \
  --tls-cert /etc/ssl/vaultpack.crt \
  --tls-key  /etc/ssl/vaultpack.key

# mTLS
vaultpack serve --listen :8443 \
  --tls-cert server.crt --tls-key server.key \
  --tls-client-ca client-ca.pem

# UNIX socket sidecar
vaultpack serve --listen unix:/run/vaultpack.sock --auth-disabled
```

## Auth requirements

The server refuses to start unless one of:

- `--auth-token` / `--auth-token-file` (constant-time bearer check)
- `--tls-client-ca` (mTLS)
- `--oidc-issuer` (+ audience / JWKS) for JWT auth
- `--auth-disabled` (explicit opt-out)

Prefer `--auth-token-file` in containers so the secret is not on argv.

`GET /metrics` requires auth whenever authentication is configured. Pass `--metrics-public` only if an unauthenticated Prometheus scraper must reach the port.

## Policy, audit, rate limits

```bash
vaultpack serve --listen :8443 \
  --auth-token-file /run/secrets/token \
  --policy /etc/vaultpack/policy.yaml \
  --audit-log /var/log/vaultpack-audit.jsonl \
  --rate-limit-rps 50 --rate-limit-burst 100
```

Serve **fails closed** if `--policy` cannot be loaded.

## KMS unwrap cache

Cached by `(provider, key-id, wrapped-DEK)` for `--kms-cache-ttl` (default 5m), max `--kms-cache-size` (default 1024). Metrics: `vaultpack_kms_cache_*`.

See [deployment.md](deployment.md) for Docker/sidecar patterns.
