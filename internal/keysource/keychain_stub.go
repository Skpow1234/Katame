//go:build !darwin

package keysource

import "fmt"

const (
	keychainAvailable = false
	keychainDetails   = "macOS-only; rebuild on darwin or use file:// / env://"
)

func resolveKeychainSymmetric(value string) ([]byte, error) {
	return nil, fmt.Errorf("keychain:// is only available on macOS (got %q)", value)
}

func resolveKeychainPrivate(value string) ([]byte, error) {
	return nil, fmt.Errorf("keychain:// is only available on macOS (got %q)", value)
}
