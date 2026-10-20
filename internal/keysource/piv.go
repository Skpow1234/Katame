package keysource

import (
	"fmt"
	"os"
	"strings"
)

const (
	pivAvailable = true
	pivDetails   = "file-backed PIV shim via VAULTPACK_PIV_DIR (slot → .pem); YubiKey piv-go in later builds"
)

// piv://slot  — e.g. piv://9a or piv://9c
//
// Reads VAULTPACK_PIV_DIR/<slot>.pem for private key material. Native YubiKey
// PIV (go-piv) is planned behind optional CGO builds.
func resolvePIVPrivate(value string) ([]byte, error) {
	slot := strings.TrimPrefix(strings.ToLower(value), "/")
	if slot == "" {
		return "", fmt.Errorf("piv URI must be piv://SLOT (e.g. piv://9a)")
	}
	dir := os.Getenv("VAULTPACK_PIV_DIR")
	if dir == "" {
		return nil, fmt.Errorf("piv://: set VAULTPACK_PIV_DIR to a directory of <slot>.pem files (native PIV support planned)")
	}
	path := strings.TrimRight(dir, `/\`) + "/" + slot + ".pem"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("piv://%s: %w", slot, err)
	}
	return data, nil
}

func resolvePIVSymmetric(value string) ([]byte, error) {
	return nil, fmt.Errorf("piv:// does not provide symmetric DEKs (use for signing keys only); got %q", value)
}
