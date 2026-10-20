// Package katame is the public, semver-stable Go SDK for Katame.
//
// It exposes the core operations of the katame CLI — Protect, Decrypt,
// Sign, Verify, Inspect — as a programmable Go API. The SDK is a thin layer
// over Katame's internal building blocks; it produces byte-for-byte
// identical .katame bundles to the CLI for matching inputs.
//
// # Stability
//
// The package follows Semantic Versioning 2.0.0. The public surface is:
//
//   - All exported identifiers in `pkg/katame`.
//   - The on-disk .katame bundle format (a ZIP archive containing
//     payload.bin, manifest.json, and optionally signature.sig and
//     provenance.json).
//   - The manifest JSON schema as documented in docs/design.md.
//
// Backwards-compatible additions (new fields with `omitempty`, new option
// values that default to current behavior) are allowed in minor releases.
// Breaking changes — renamed fields, removed functions, changed semantics —
// require a major version bump.
//
// # Quick start
//
//	import "github.com/Skpow1234/Katame/pkg/katame"
//
//	// Encrypt a file with a freshly-generated key.
//	res, err := katame.Protect(katame.ProtectOptions{
//	    InputPath:  "secret.txt",
//	    OutputPath: "secret.katame",
//	})
//	if err != nil { panic(err) }
//	os.WriteFile("secret.key", res.GeneratedKey, 0o600)
//
//	// Decrypt it later.
//	dec, err := katame.Decrypt(katame.DecryptOptions{
//	    InputPath:  "secret.katame",
//	    OutputPath: "secret.out",
//	    Key:        keyBytes,
//	})
//	if err != nil { panic(err) }
//	_ = dec
//
// # Supported modes (v1.0)
//
//   - Symmetric key file (generated or supplied).
//   - Password-based (PBKDF: argon2id, scrypt, pbkdf2-sha256).
//   - Hybrid / multi-recipient encryption (PEM public keys).
//   - KMS-wrapped DEKs (aws, gcp, azure, mock providers).
//   - Pre-encryption compression (gzip, zstd).
//   - Shamir key splitting (K-of-N shares via SplitShares / CombineShares).
//   - Detached signatures (ed25519 / ecdsa / rsa-pss / ml-dsa / slh-dsa).
//   - Sigstore Rekor transparency upload and verification helpers.
//   - SLSA-style provenance embedding (provenance.json).
//   - Bundle inspection (manifest decoding).
package katame
