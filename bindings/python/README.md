# katame — Python bindings

`pip install katame` provides a small `ctypes` wrapper around
`libkatame`, the Katame C-shared library. No C compiler is required
at install time, but the native library must be present on the host.

## Installation

```bash
pip install katame
```

Then make `libkatame.{so,dylib,dll}` discoverable by one of:

- Setting `$KATAME_LIB=/abs/path/to/libkatame.so`.
- Dropping it next to the installed `katame/` package directory.
- Putting it on `LD_LIBRARY_PATH` / `DYLD_LIBRARY_PATH` / `PATH` so the
  system loader finds it.

Build the library from the Katame source tree:

```bash
go build -buildmode=c-shared -tags cshared \
  -o libkatame.so ./cmd/katame-c
```

(See `cmd/katame-c/README.md` for per-OS instructions.)

## Quick start

```python
import base64
import katame

# Encrypt with a generated key.
res = katame.protect(
    input_path="secret.txt",
    output_path="secret.katame",
)
key_b64 = res["generated_key_b64"]   # store this somewhere safe!
print("manifest:", res["manifest"])

# Decrypt back.
out = katame.decrypt(
    input_path="secret.katame",
    key_b64=key_b64,
    output_path="secret.out",
)
print("decrypted manifest:", out["manifest"])

# Or do it entirely in memory:
key = base64.b64decode(key_b64)
res2 = katame.decrypt(
    input_path="secret.katame",
    key=key,
)
plaintext = base64.b64decode(res2["plaintext_b64"])
print(plaintext)
```

## API

| Function                          | Purpose                              |
|-----------------------------------|--------------------------------------|
| `version() -> str`                | libkatame semver                  |
| `protect(...) -> dict`            | Encrypt → .katame bundle              |
| `decrypt(...) -> dict`            | Decrypt a .katame bundle              |
| `inspect(path) -> dict`           | Read the manifest only               |
| `sign(...) -> dict`               | Add/replace detached signature       |
| `verify(...) -> dict`             | Verify detached signature            |

Every call returns a `dict` (JSON-deserialized). On error a
`katame.KatameError` is raised with the underlying message.

## Supported modes (v1.0)

- Symmetric key file (auto-generated or supplied).
- Password-based (argon2id / scrypt / pbkdf2-sha256).
- Hybrid recipients (`recipients=[{"public_key_pem": "..."}]`).
- KMS wrap (`kms_provider` / `kms_key_id`).
- Compression and Shamir split (`compress`, `split_shares` / `split_threshold`).
- Detached signatures (ed25519 / ecdsa / rsa-pss / ml-dsa / slh-dsa).
- KMS rewrap via `rewrap(...)`.

Sigstore transparency uploads remain CLI-oriented for now.

## License

MIT.
