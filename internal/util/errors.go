package util

import "errors"

// Exit codes for automation-friendly CLI usage.
const (
	ExitSuccess            = 0
	ExitGenericError       = 1
	ExitInvalidArgs        = 2
	ExitVerifyFailed       = 10
	ExitDecryptFailed      = 11
	ExitUnsupportedVersion = 12
	ExitPolicyDenied       = 13
)

// Sentinel errors used across the application.
var (
	ErrUnsupportedVersion   = errors.New("unsupported manifest version")
	ErrUnsupportedAlgorithm = errors.New("unsupported algorithm")
	ErrInvalidNonceLength   = errors.New("invalid nonce length")
	ErrInvalidTagLength     = errors.New("invalid tag length")
	ErrInvalidKeyLength     = errors.New("invalid key length")
	ErrDecryptFailed        = errors.New("decryption failed")
	ErrVerifyFailed         = errors.New("signature verification failed")
	ErrManifestInvalid      = errors.New("invalid manifest")
	ErrBundleCorrupted      = errors.New("bundle is corrupted or incomplete")
	ErrKeyMismatch          = errors.New("key fingerprint does not match manifest")
	ErrPolicyDenied         = errors.New("policy denied")
	ErrBundlesDiffer        = errors.New("bundles differ")
)

// ExitCodeForError maps a sentinel error to its CLI exit code.
func ExitCodeForError(err error) int {
	switch {
	case err == nil:
		return ExitSuccess
	case errors.Is(err, ErrPolicyDenied):
		return ExitPolicyDenied
	case errors.Is(err, ErrVerifyFailed), errors.Is(err, ErrBundlesDiffer):
		return ExitVerifyFailed
	case errors.Is(err, ErrDecryptFailed), errors.Is(err, ErrKeyMismatch), errors.Is(err, ErrBundleCorrupted):
		return ExitDecryptFailed
	case errors.Is(err, ErrUnsupportedVersion), errors.Is(err, ErrUnsupportedAlgorithm):
		return ExitUnsupportedVersion
	default:
		return ExitGenericError
	}
}
