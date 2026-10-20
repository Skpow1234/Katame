# Security

## Algorithms

- **AEAD:** AES-256-GCM, ChaCha20-Poly1305, XChaCha20-Poly1305, AES-256-GCM-SIV
- **Hashing:** SHA-256/512, SHA3-256/512, BLAKE2b-256/512, BLAKE3
- **Signing:** Ed25519, Ed448, ECDSA P-256/P-384, RSA-PSS 2048/4096, ML-DSA-65/87, SLH-DSA-128s
- **KDFs:** Argon2id (t=3, m=64 MiB, p=4), scrypt (N=32768), PBKDF2-SHA256 (600k)
- **Hybrid:** X25519+HKDF+AES-256-GCM, ECIES-P256, RSA-OAEP-SHA256, ML-KEM-768/1024

## Chunked streaming encryption

Default 64 KiB plaintext chunks. VaultPack mitigates the three standard chunked-AEAD attacks:

**Nonce derivation.** Random base nonce per bundle; chunk nonce = `base_nonce XOR chunk_index` (big-endian 64-bit counter in the last 8 bytes).

**Last-chunk flag.** Final chunk sets bit 63 on the counter (`counter | 0x8000000000000000`) so truncation fails AEAD auth.

**Ordering.** Nonce encodes sequential index — reorder/duplicate/remove fails authentication.

Per chunk: plaintext + nonce + optional user AAD. Across the bundle: manifest records base nonce, final tag, size, cipher, chunk size; signatures cover canonical manifest + SHA-256 of `payload.bin`.

## Other properties

- Key fingerprint (SHA-256 of raw key) for early wrong-key detection
- KMS: DEK wrapped; only wrapped ciphertext + key ID stored
- Multi-recipient: independent wraps of one DEK
- ECDH hybrid: ephemeral keys → forward secrecy for past DEKs
- Compression runs *before* encryption (no compression-ratio oracle on ciphertext)
- Shamir GF(256): fewer than K shares reveal zero information
- Passwords/keys/private keys are never stored inside the bundle

## Operational notes

- Prefer signed bundles + policy deny for unsigned decrypts in production
- Do not expose `vaultpack serve --auth-disabled` on TCP
- Treat `/metrics` as sensitive (unauthenticated by default)
- Plugin binaries are a trust boundary (DEK material via env/JSON)
