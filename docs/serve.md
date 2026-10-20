# Service mode (`katame serve`)

Runs a long-running HTTP API exposing the SDK. The wire contract matches the C-shared library and WASM bindings: JSON in/out with an `{"ok": true|false, ...}` envelope.

Full route reference: [endpoints.md](endpoints.md) · OpenAPI: [openapi.yaml](openapi.yaml).

## Examples

```bash
# Local development only
katame serve --listen :8443 --auth-disabled

# Bearer token + TLS
katame serve \
  --listen :8443 \
  --auth-token "$(cat /etc/katame/token)" \
  --tls-cert /etc/ssl/katame.crt \
  --tls-key  /etc/ssl/katame.key

# mTLS
katame serve --listen :8443 \
  --tls-cert server.crt --tls-key server.key \
  --tls-client-ca client-ca.pem

# UNIX socket sidecar
katame serve --listen unix:/run/katame.sock --auth-disabled
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
katame serve --listen :8443 \
  --auth-token-file /run/secrets/token \
  --policy /etc/katame/policy.yaml \
  --audit-log /var/log/katame-audit.jsonl \
  --rate-limit-rps 50 --rate-limit-burst 100
```

Serve **fails closed** if `--policy` cannot be loaded.

## KMS unwrap cache

Cached by `(provider, key-id, wrapped-DEK)` for `--kms-cache-ttl` (default 5m), max `--kms-cache-size` (default 1024). Metrics: `katame_kms_cache_*`.

See [deployment.md](deployment.md) for Docker/sidecar patterns.
