//go:build !(cgo && katame_pkcs11)

package keysource

import (
	stdcrypto "crypto"
	"fmt"
)

const nativePKCS11Enabled = false

func resolveNativePKCS11Symmetric(ref pkcs11Ref) ([]byte, error) {
	return nil, fmt.Errorf("pkcs11://%s: native SoftHSM not in this binary (rebuild with CGO_ENABLED=1 -tags katame_pkcs11)", ref.Label)
}

func resolveNativePKCS11Signer(ref pkcs11Ref) (stdcrypto.Signer, string, error) {
	return nil, "", fmt.Errorf("pkcs11://%s: native SoftHSM not in this binary (rebuild with CGO_ENABLED=1 -tags katame_pkcs11)", ref.Label)
}
