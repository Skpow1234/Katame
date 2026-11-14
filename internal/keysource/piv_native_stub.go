//go:build !(cgo && katame_piv)

package keysource

import (
	stdcrypto "crypto"
	"fmt"
)

const nativePIVEnabled = false

func resolveNativePIVSigner(ref pivRef) (stdcrypto.Signer, string, error) {
	return nil, "", fmt.Errorf("piv://%s: native YubiKey PIV not in this binary (rebuild with CGO_ENABLED=1 -tags katame_piv)", ref.Slot)
}
