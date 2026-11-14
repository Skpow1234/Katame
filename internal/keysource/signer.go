package keysource

import (
	stdcrypto "crypto"
	"fmt"

	vpcrypto "github.com/Skpow1234/Katame/internal/crypto"
)

// ResolveSigner loads a crypto.Signer from a key-source URI.
//
// Portable schemes (file/env/b64/dpapi/keychain) and the file-backed
// pkcs11:// / piv:// shims resolve key bytes then parse them. Native
// SoftHSM (M251, -tags katame_pkcs11) and YubiKey PIV (M252, -tags
// katame_piv) backends return hardware-backed signers that never export
// private key material.
func ResolveSigner(uri string) (stdcrypto.Signer, string, error) {
	scheme, value, err := split(uri)
	if err != nil {
		return nil, "", err
	}
	switch scheme {
	case "piv":
		return resolvePIVSigner(value)
	case "pkcs11":
		return resolvePKCS11Signer(value)
	default:
		keyBytes, err := ResolvePrivateKey(uri)
		if err != nil {
			return nil, "", err
		}
		return vpcrypto.ParsePrivateKey(keyBytes)
	}
}

func signerFromPrivateBytes(data []byte) (stdcrypto.Signer, string, error) {
	return vpcrypto.ParsePrivateKey(data)
}

func algoFromSigner(signer stdcrypto.Signer) (string, error) {
	if signer == nil {
		return "", fmt.Errorf("nil signer")
	}
	return vpcrypto.AlgoFromPublicKey(signer.Public())
}
