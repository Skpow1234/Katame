//go:build darwin

package keysource

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const (
	keychainAvailable = true
	keychainDetails   = "macOS Keychain via /usr/bin/security (generic password)"
)

// keychain://SERVICE/ACCOUNT — reads a generic-password secret (-w).
// The secret must be a Katame key file body (b64:…) or PEM private key.
func resolveKeychainSymmetric(value string) ([]byte, error) {
	raw, err := keychainFindPassword(value)
	if err != nil {
		return nil, err
	}
	return parseSymmetricKeyBytes([]byte(raw))
}

func resolveKeychainPrivate(value string) ([]byte, error) {
	raw, err := keychainFindPassword(value)
	if err != nil {
		return nil, err
	}
	return []byte(raw), nil
}

func keychainFindPassword(value string) (string, error) {
	service, account, err := splitKeychainRef(value)
	if err != nil {
		return "", err
	}
	cmd := exec.Command("/usr/bin/security", "find-generic-password", "-s", service, "-a", account, "-w")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("keychain://%s/%s: %w (%s)", service, account, err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}

func splitKeychainRef(value string) (service, account string, err error) {
	value = strings.TrimPrefix(value, "/")
	parts := strings.SplitN(value, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("keychain URI must be keychain://SERVICE/ACCOUNT, got %q", value)
	}
	return parts[0], parts[1], nil
}
