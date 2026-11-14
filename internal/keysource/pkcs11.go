package keysource

import (
	stdcrypto "crypto"
	"fmt"
	"os"
	"strings"
)

const (
	pkcs11Available = true
)

var pkcs11Details = pkcs11DetailsString()

func pkcs11DetailsString() string {
	if nativePKCS11Enabled {
		return "native PKCS#11 via KATAME_PKCS11_MODULE (SoftHSM); file shim via KATAME_PKCS11_DIR"
	}
	return "file-backed PKCS#11 shim via KATAME_PKCS11_DIR; rebuild with CGO_ENABLED=1 -tags katame_pkcs11 for SoftHSM"
}

// pkcs11://label[?pin=…&slot=N]
//
// When built with -tags katame_pkcs11 and KATAME_PKCS11_MODULE is set, keys are
// loaded from the PKCS#11 token (SoftHSM). Otherwise labels resolve from
// KATAME_PKCS11_DIR/<label>.key (symmetric) or .pem (private).
func resolvePKCS11Symmetric(value string) ([]byte, error) {
	ref, err := parsePKCS11Ref(value)
	if err != nil {
		return nil, err
	}
	if nativePKCS11Enabled && pkcs11ModulePath() != "" {
		key, err := resolveNativePKCS11Symmetric(ref)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	return resolvePKCS11SymmetricShim(ref)
}

func resolvePKCS11Private(value string) ([]byte, error) {
	ref, err := parsePKCS11Ref(value)
	if err != nil {
		return nil, err
	}
	if nativePKCS11Enabled && pkcs11ModulePath() != "" {
		return nil, fmt.Errorf("pkcs11://%s: native SoftHSM private keys are non-extractable; use ResolveSigner / --signing-key-source with a crypto.Signer backend", ref.Label)
	}
	return resolvePKCS11PrivateShim(ref)
}

func resolvePKCS11Signer(value string) (stdcrypto.Signer, string, error) {
	ref, err := parsePKCS11Ref(value)
	if err != nil {
		return nil, "", err
	}
	if nativePKCS11Enabled && pkcs11ModulePath() != "" {
		return resolveNativePKCS11Signer(ref)
	}
	pemBytes, err := resolvePKCS11PrivateShim(ref)
	if err != nil {
		return nil, "", err
	}
	return signerFromPrivateBytes(pemBytes)
}

func resolvePKCS11SymmetricShim(ref pkcs11Ref) ([]byte, error) {
	dir := os.Getenv("KATAME_PKCS11_DIR")
	if dir == "" {
		if nativePKCS11Enabled {
			return nil, fmt.Errorf("pkcs11://: set KATAME_PKCS11_MODULE for SoftHSM or KATAME_PKCS11_DIR for the file shim")
		}
		return nil, fmt.Errorf("pkcs11://: set KATAME_PKCS11_DIR to a directory of <label>.key files (or rebuild with -tags katame_pkcs11)")
	}
	path := strings.TrimRight(dir, `/\`) + "/" + ref.Label + ".key"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("pkcs11://%s: %w", ref.Label, err)
	}
	return parseSymmetricKeyBytes(data)
}

func resolvePKCS11PrivateShim(ref pkcs11Ref) ([]byte, error) {
	dir := os.Getenv("KATAME_PKCS11_DIR")
	if dir == "" {
		if nativePKCS11Enabled {
			return nil, fmt.Errorf("pkcs11://: set KATAME_PKCS11_MODULE for SoftHSM or KATAME_PKCS11_DIR for the file shim")
		}
		return nil, fmt.Errorf("pkcs11://: set KATAME_PKCS11_DIR to a directory of <label>.pem files (or rebuild with -tags katame_pkcs11)")
	}
	path := strings.TrimRight(dir, `/\`) + "/" + ref.Label + ".pem"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("pkcs11://%s: %w", ref.Label, err)
	}
	return data, nil
}

type pkcs11Ref struct {
	Label string
	PIN   string
	Slot  int // -1 = first token slot
}

func parsePKCS11Ref(value string) (pkcs11Ref, error) {
	ref := pkcs11Ref{Slot: -1}
	value = strings.TrimPrefix(value, "/")
	if value == "" {
		return ref, fmt.Errorf("pkcs11 URI must include a label")
	}
	base, query, hasQuery := strings.Cut(value, "?")
	ref.Label = base
	if i := strings.LastIndex(base, "/"); i >= 0 {
		ref.Label = base[i+1:]
	}
	if ref.Label == "" {
		return ref, fmt.Errorf("pkcs11 URI missing label")
	}
	if hasQuery {
		for _, part := range strings.Split(query, "&") {
			k, v, ok := strings.Cut(part, "=")
			if !ok {
				continue
			}
			switch k {
			case "pin":
				ref.PIN = v
			case "label":
				if v != "" {
					ref.Label = v
				}
			case "slot":
				var n int
				if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
					ref.Slot = n
				}
			}
		}
	}
	if ref.PIN == "" {
		ref.PIN = os.Getenv("KATAME_PKCS11_PIN")
	}
	return ref, nil
}

func pkcs11ModulePath() string {
	if p := os.Getenv("KATAME_PKCS11_MODULE"); p != "" {
		return p
	}
	// Common SoftHSM2 install locations (used when building with katame_pkcs11).
	for _, cand := range []string{
		"/usr/lib/softhsm/libsofthsm2.so",
		"/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so",
		"/usr/local/lib/softhsm/libsofthsm2.so",
		"/opt/homebrew/lib/softhsm/libsofthsm2.so",
	} {
		if st, err := os.Stat(cand); err == nil && !st.IsDir() {
			return cand
		}
	}
	return ""
}
