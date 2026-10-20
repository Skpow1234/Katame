package keysource

import (
	"fmt"
	"os"
	"strings"
)

const (
	pkcs11Available = true
	pkcs11Details   = "file-backed PKCS#11 shim via VAULTPACK_PKCS11_DIR (label → .key/.pem); native SoftHSM in later builds"
)

// pkcs11://label[?pin=…] or pkcs11://module/path?label=NAME
//
// Production HSMs need a native PKCS#11 module. This build resolves labels from
// VAULTPACK_PKCS11_DIR/<label>.key (symmetric) or .pem (private) so CI and
// doctor can exercise the URI scheme without CGO SoftHSM bindings.
func resolvePKCS11Symmetric(value string) ([]byte, error) {
	label, _, err := parsePKCS11Ref(value)
	if err != nil {
		return nil, err
	}
	dir := os.Getenv("VAULTPACK_PKCS11_DIR")
	if dir == "" {
		return nil, fmt.Errorf("pkcs11://: set VAULTPACK_PKCS11_DIR to a directory of <label>.key files (native module support planned)")
	}
	path := strings.TrimRight(dir, `/\`) + "/" + label + ".key"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("pkcs11://%s: %w", label, err)
	}
	return parseSymmetricKeyBytes(data)
}

func resolvePKCS11Private(value string) ([]byte, error) {
	label, _, err := parsePKCS11Ref(value)
	if err != nil {
		return nil, err
	}
	dir := os.Getenv("VAULTPACK_PKCS11_DIR")
	if dir == "" {
		return nil, fmt.Errorf("pkcs11://: set VAULTPACK_PKCS11_DIR to a directory of <label>.pem files (native module support planned)")
	}
	path := strings.TrimRight(dir, `/\`) + "/" + label + ".pem"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("pkcs11://%s: %w", label, err)
	}
	return data, nil
}

func parsePKCS11Ref(value string) (label, pin string, err error) {
	value = strings.TrimPrefix(value, "/")
	if value == "" {
		return "", "", fmt.Errorf("pkcs11 URI must include a label")
	}
	base, query, hasQuery := strings.Cut(value, "?")
	label = base
	if i := strings.LastIndex(base, "/"); i >= 0 {
		label = base[i+1:]
	}
	if label == "" {
		return "", "", fmt.Errorf("pkcs11 URI missing label")
	}
	if hasQuery {
		for _, part := range strings.Split(query, "&") {
			k, v, ok := strings.Cut(part, "=")
			if ok && (k == "pin" || k == "label") {
				if k == "pin" {
					pin = v
				}
				if k == "label" && v != "" {
					label = v
				}
			}
		}
	}
	return label, pin, nil
}
