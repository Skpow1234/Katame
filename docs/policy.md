# Policy & RBAC

Policy is resolved: `--policy <path>` → `VAULTPACK_POLICY` → `policy_file:` in `.vpack.yaml`. With no policy configured, operations run unrestricted.

**Fail-closed:** If `--policy` / `VAULTPACK_POLICY` / `policy_file` is set and the file cannot be loaded, the CLI aborts (same as `vaultpack serve`). Validate with `vaultpack policy validate` before deploying.

## YAML example

```yaml
version: 1
default: allow
rules:
  - name: block-dev-keys-in-prod
    action: deny
    reason: "dev KMS keys cannot be decrypted from prod"
    when:
      operation: [decrypt, rewrap]
      kms_key_id_matches: ".*/dev-.*"

  - name: deny-unsigned-decrypts
    action: deny
    reason: "decrypts require a signature"
    when:
      operation: [decrypt]
      is_signed: false

  - name: after-hours
    action: deny
    reason: "no decrypts between 00:00 and 06:00 UTC"
    when:
      operation: [decrypt]
      time_window: {from: "00:00", to: "06:00", timezone: "UTC"}
```

Rules are first-match. Predicates within a rule are ANDed.

| Field | Description |
| ----- | ----------- |
| `operation` | `decrypt`, `verify`, `protect`, … |
| `user` / `user_matches` | Exact or regex username |
| `hostname` / `hostname_matches` | Host |
| `weekday` | Day names |
| `time_window` | `{from, to, timezone}` HH:MM |
| `bundle_path_matches` | Regex on path |
| `kms_key_id` / `kms_key_id_matches` | KMS key |
| `cipher_in` | Allowed AEAD names |
| `hybrid_scheme_matches` | Regex on hybrid scheme |
| `signature_algo_in` / `signature_algo_not_in` | Signature algos |
| `is_signed` | Bool |
| `recipient_fingerprint_in` | Fingerprints |
| `plaintext_digest_in` | Plaintext digests (base64) |
| `min_recipients` / `max_recipients` | Bounds |

JSON files use the same field names.

## Rego (OPA)

`.rego` files load automatically. Expose `allow` and optional `deny_reason`:

```rego
package vaultpack

default allow = false

allow {
    input.operation == "decrypt"
    input.user == "alice"
    startswith(input.manifest.kms_key_id, "arn:aws:kms:us-east-1:")
}

deny_reason = sprintf("user %q cannot decrypt %q", [input.user, input.manifest.kms_key_id]) {
    not allow
}
```

## Commands

```bash
vaultpack policy validate --file policy.yaml
vaultpack policy test --file policy.yaml --bundle secret.vpack --op decrypt --user alice
vaultpack policy show --policy policy.yaml
```

Denials are audited as `policy-deny` with rule name and reason.
