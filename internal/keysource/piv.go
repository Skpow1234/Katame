package keysource

import (
	stdcrypto "crypto"
	"fmt"
	"os"
	"strings"
)

const (
	pivAvailable = true
)

var pivDetails = pivDetailsString()

func pivDetailsString() string {
	if nativePIVEnabled {
		return "native YubiKey PIV via go-piv; file shim via KATAME_PIV_DIR"
	}
	return "file-backed PIV shim via KATAME_PIV_DIR; rebuild with CGO_ENABLED=1 -tags katame_piv for YubiKey"
}

// piv://slot[?pin=…]  — e.g. piv://9a or piv://9c?pin=123456
//
// File shim reads KATAME_PIV_DIR/<slot>.pem. Native builds (-tags katame_piv)
// open a YubiKey via go-piv and return a hardware crypto.Signer.
func resolvePIVPrivate(value string) ([]byte, error) {
	ref, err := parsePIVRef(value)
	if err != nil {
		return nil, err
	}
	if nativePIVEnabled {
		return nil, fmt.Errorf("piv://%s: YubiKey private keys are non-extractable; use ResolveSigner / --signing-key-source", ref.Slot)
	}
	return resolvePIVPrivateShim(ref)
}

func resolvePIVSymmetric(value string) ([]byte, error) {
	return nil, fmt.Errorf("piv:// does not provide symmetric DEKs (use for signing keys only); got %q", value)
}

func resolvePIVSigner(value string) (stdcrypto.Signer, string, error) {
	ref, err := parsePIVRef(value)
	if err != nil {
		return nil, "", err
	}
	if nativePIVEnabled {
		return resolveNativePIVSigner(ref)
	}
	pemBytes, err := resolvePIVPrivateShim(ref)
	if err != nil {
		return nil, "", err
	}
	return signerFromPrivateBytes(pemBytes)
}

func resolvePIVPrivateShim(ref pivRef) ([]byte, error) {
	dir := os.Getenv("KATAME_PIV_DIR")
	if dir == "" {
		if nativePIVEnabled {
			return nil, fmt.Errorf("piv://: insert a YubiKey or set KATAME_PIV_DIR for the file shim")
		}
		return nil, fmt.Errorf("piv://: set KATAME_PIV_DIR to a directory of <slot>.pem files (or rebuild with -tags katame_piv)")
	}
	path := strings.TrimRight(dir, `/\`) + "/" + ref.Slot + ".pem"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("piv://%s: %w", ref.Slot, err)
	}
	return data, nil
}

type pivRef struct {
	Slot string
	PIN  string
}

func parsePIVRef(value string) (pivRef, error) {
	ref := pivRef{}
	value = strings.TrimPrefix(value, "/")
	base, query, hasQuery := strings.Cut(value, "?")
	ref.Slot = strings.ToLower(base)
	if ref.Slot == "" {
		return ref, fmt.Errorf("piv URI must be piv://SLOT (e.g. piv://9a)")
	}
	if hasQuery {
		for _, part := range strings.Split(query, "&") {
			k, v, ok := strings.Cut(part, "=")
			if ok && k == "pin" {
				ref.PIN = v
			}
		}
	}
	if ref.PIN == "" {
		ref.PIN = os.Getenv("KATAME_PIV_PIN")
	}
	return ref, nil
}
