# HTTP endpoints

Base: `vaultpack serve`. All `/v1/*` replies use `{"ok": true|false, ...}`.

| Method | Path | Auth | Purpose |
| ------ | ---- | ---- | ------- |
| `GET` | `/healthz` | no | Liveness |
| `GET` | `/metrics` | no | Prometheus text exposition |
| `GET` | `/v1/version` | yes | SDK + API version |
| `POST` | `/v1/protect` | yes | Encrypt → bundle |
| `POST` | `/v1/decrypt` | yes | Decrypt bundle |
| `POST` | `/v1/inspect` | yes | Decode manifest |
| `POST` | `/v1/sign` | yes | Add/replace detached signature |
| `POST` | `/v1/verify` | yes | Verify signature |
| `POST` | `/v1/rewrap` | yes | Re-wrap KMS DEK |

Machine-readable contract: [openapi.yaml](openapi.yaml).

## Authentication

Send `Authorization: Bearer <token>` when using bearer auth, or present a client certificate for mTLS / OIDC JWT as configured.

`/healthz` and `/metrics` skip auth by design — restrict network exposure or put metrics behind a sidecar proxy in production.

## curl round-trip

```bash
TOKEN="$(cat /etc/vaultpack/token)"

RESPONSE=$(curl -s -X POST https://vp.example.com/v1/protect \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"plaintext_b64":"aGVsbG8gd29ybGQ="}')
BUNDLE=$(jq -r .bundle_b64 <<< "$RESPONSE")
KEY=$(jq -r .generated_key_b64 <<< "$RESPONSE")

curl -s -X POST https://vp.example.com/v1/decrypt \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "$(jq -n --arg b "$BUNDLE" --arg k "$KEY" '{bundle_b64:$b, key_b64:$k}')" \
  | jq -r .plaintext_b64 | base64 -d
```

## Request shapes (summary)

### `POST /v1/protect`

Common fields: `plaintext_b64` or `input_path`, optional `output_path`, `key_b64`, `password`, `kdf_algo`, `kms_provider`, `kms_key_id`, `recipients[]`, `compress`, `split_shares` / `split_threshold`, `cipher`, `hash_algo`, `aad_b64`, signing fields, `return_bundle`.

Response may include `bundle_b64`, `generated_key_b64`, `shares`, `manifest`, `signature_*`.

### `POST /v1/decrypt`

`bundle_b64` or `input_path`, plus one of `key_b64`, `password`, `private_key_pem` / `private_key_path`, or `kms_provider`. Optional `return_plaintext`.

### `POST /v1/inspect`

`input_path` or `bundle_b64` → `{ok, manifest}`.

### `POST /v1/sign` / `POST /v1/verify`

Bundle path or bytes + key material / pubkey.

### `POST /v1/rewrap`

KMS rewrap of an existing wrapped DEK (rotation without touching payload ciphertext).

## Metrics

Prometheus 0.0.4 text at `GET /metrics`:

- `vaultpack_build_info{version}`
- `vaultpack_uptime_seconds`
- `vaultpack_http_requests_total{path,status}`
- `vaultpack_http_request_duration_seconds` (histogram)
- `vaultpack_auth_denied_total`
- `vaultpack_kms_cache_entries` / `_hits_total` / `_misses_total` / `_evicted_total`

Default body limit: 64 MiB (`--max-request-bytes`).
