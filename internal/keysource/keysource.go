// Package keysource resolves key material from URI-style references.
//
// M28 establishes the stable abstraction used by CLI flags such as
// --key-source and --signing-key-source. Portable backends (file/env/b64)
// plus platform schemes (dpapi, keychain, pkcs11, piv) are implemented here.
package keysource

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	vpcrypto "github.com/Skpow1234/Vaultpack/internal/crypto"
	"github.com/Skpow1234/Vaultpack/internal/util"
)

// UnsupportedSchemeError is returned for unknown schemes.
type UnsupportedSchemeError struct {
	Scheme string
	URI    string
}

func (e UnsupportedSchemeError) Error() string {
	return fmt.Sprintf("key source scheme %q is not supported: %s", e.Scheme, e.URI)
}

// ResolveSymmetricKey loads a raw 32-byte VaultPack content-encryption key.
//
// Supported forms:
//   - file:///abs/path.key
//   - file://relative/or/windows/path.key
//   - env://ENV_VAR_CONTAINING_KEYFILE_OR_BASE64
//   - b64://BASE64_RAW_32_BYTE_KEY
//   - dpapi://path|b64:… (Windows)
//   - keychain://SERVICE/ACCOUNT (macOS)
//   - pkcs11://label (VAULTPACK_PKCS11_DIR)
func ResolveSymmetricKey(uri string) ([]byte, error) {
	scheme, value, err := split(uri)
	if err != nil {
		return nil, err
	}
	switch scheme {
	case "file":
		return vpcrypto.LoadKeyFile(value)
	case "env":
		raw, ok := os.LookupEnv(value)
		if !ok {
			return nil, fmt.Errorf("env key source %q is not set", value)
		}
		return parseSymmetricKeyBytes([]byte(raw))
	case "b64":
		return parseSymmetricKeyBytes([]byte(value))
	case "dpapi":
		return resolveDPAPISymmetric(value)
	case "keychain":
		return resolveKeychainSymmetric(value)
	case "pkcs11":
		return resolvePKCS11Symmetric(value)
	case "piv":
		return resolvePIVSymmetric(value)
	default:
		return nil, UnsupportedSchemeError{Scheme: scheme, URI: uri}
	}
}

// ResolvePrivateKey loads a signing/private key as bytes. Callers should pass
// the result to crypto.ParsePrivateKey so all existing key formats remain
// supported.
//
// Supported forms:
//   - file:///abs/signing.key
//   - env://ENV_VAR_CONTAINING_PEM_OR_LEGACY_KEY
//   - b64://BASE64_PEM_OR_LEGACY_KEY_BYTES
//   - dpapi://path|b64:… (Windows)
//   - keychain://SERVICE/ACCOUNT (macOS)
//   - pkcs11://label (VAULTPACK_PKCS11_DIR)
//   - piv://slot (VAULTPACK_PIV_DIR)
func ResolvePrivateKey(uri string) ([]byte, error) {
	scheme, value, err := split(uri)
	if err != nil {
		return nil, err
	}
	switch scheme {
	case "file":
		return os.ReadFile(value)
	case "env":
		raw, ok := os.LookupEnv(value)
		if !ok {
			return nil, fmt.Errorf("env key source %q is not set", value)
		}
		return []byte(raw), nil
	case "b64":
		return util.B64Decode(value)
	case "dpapi":
		return resolveDPAPIPrivate(value)
	case "keychain":
		return resolveKeychainPrivate(value)
	case "pkcs11":
		return resolvePKCS11Private(value)
	case "piv":
		return resolvePIVPrivate(value)
	default:
		return nil, UnsupportedSchemeError{Scheme: scheme, URI: uri}
	}
}

func split(uri string) (scheme, value string, err error) {
	if uri == "" {
		return "", "", fmt.Errorf("key source URI is empty")
	}
	i := strings.Index(uri, "://")
	if i <= 0 {
		return "", "", fmt.Errorf("key source %q must use scheme://value syntax", uri)
	}
	scheme = strings.ToLower(uri[:i])
	rawValue := uri[i+3:]
	if rawValue == "" {
		return "", "", fmt.Errorf("key source %q has empty value", uri)
	}
	if scheme == "file" {
		return scheme, fileValue(rawValue), nil
	}
	return scheme, rawValue, nil
}

func fileValue(raw string) string {
	// url.PathUnescape keeps Windows paths like C:/... intact while allowing
	// spaces to be represented as %20.
	if decoded, err := url.PathUnescape(raw); err == nil {
		raw = decoded
	}
	if strings.HasPrefix(raw, "/") && len(raw) >= 3 && raw[2] == ':' {
		// file:///C:/Users/... -> C:/Users/... on Windows-like paths.
		return raw[1:]
	}
	return raw
}

func parseSymmetricKeyBytes(data []byte) ([]byte, error) {
	content := strings.TrimSpace(string(data))
	content = strings.TrimPrefix(content, vpcrypto.KeyFilePrefix)
	key, err := util.B64Decode(content)
	if err != nil {
		return nil, fmt.Errorf("decode symmetric key source: %w", err)
	}
	if len(key) != vpcrypto.AES256KeySize {
		return nil, fmt.Errorf("%w: got %d bytes, want %d", util.ErrInvalidKeyLength, len(key), vpcrypto.AES256KeySize)
	}
	return key, nil
}
