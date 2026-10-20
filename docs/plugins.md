# Plugin author guide

Plugins are executables discovered from `KATAME_PLUGIN_DIR` or config `plugin_dir`. Each plugin is probed with `--katame-list` and may register KEM schemes and/or signing algorithms. First plugin wins per name.

Example: [`examples/plugin-dummy-kem`](../examples/plugin-dummy-kem) (insecure demo only).

## Discovery

```bash
export KATAME_PLUGIN_DIR=/path/to/plugins
katame keygen --out demo --algo <plugin-scheme>
```

## Capability listing

When invoked as `plugin --katame-list`, stdout must be JSON:

```json
{"kem": ["my-kem"], "sign": ["my-sign"]}
```

Either array may be empty/omitted.

## KEM contract

### `keygen`

```text
plugin keygen --type kem --algo <scheme> --out <prefix>
```

Writes `<prefix>.key` and `<prefix>.pub` (may annotate pub with `KATAME-KEM-SCHEME`).

### `encapsulate`

```text
plugin encapsulate --algo <scheme> --pubkey <path>
```

Optional env `KATAME_DEK_B64` = base64 DEK to wrap (multi-recipient). Stdout JSON:

```json
{
  "dek_b64": "...",
  "ephemeral_b64": "...",
  "wrapped_dek_b64": "...",
  "scheme": "my-kem"
}
```

### `decapsulate`

```text
plugin decapsulate --algo <scheme> --privkey <path>
```

Input via `KATAME_CIPHERTEXT_B64` (base64 of JSON with `ephemeral_b64` / `wrapped_dek_b64`). Stdout: DEK as a single line of base64 or hex.

## Sign contract

Similar subprocess interface for keygen / sign / verify (see `internal/plugin/sign.go`). Registered algo names appear in `keygen --algo`, `sign`, and `verify`.

## Trust model

Plugins run as subprocesses and may receive key material via environment variables. Only load plugins from trusted directories.
