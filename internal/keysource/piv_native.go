//go:build cgo && katame_piv

package keysource

import (
	stdcrypto "crypto"
	"fmt"
	"strings"

	"github.com/go-piv/piv-go/v2/piv"
)

const nativePIVEnabled = true

func resolveNativePIVSigner(ref pivRef) (stdcrypto.Signer, string, error) {
	cards, err := piv.Cards()
	if err != nil {
		return nil, "", fmt.Errorf("piv://%s: list cards: %w", ref.Slot, err)
	}
	if len(cards) == 0 {
		return nil, "", fmt.Errorf("piv://%s: no YubiKey detected (and KATAME_PIV_DIR unset)", ref.Slot)
	}
	yk, err := piv.Open(cards[0])
	if err != nil {
		return nil, "", fmt.Errorf("piv://%s: open %q: %w", ref.Slot, cards[0], err)
	}

	slot, err := pivSlot(ref.Slot)
	if err != nil {
		yk.Close()
		return nil, "", err
	}

	cert, err := yk.Certificate(slot)
	if err != nil {
		yk.Close()
		return nil, "", fmt.Errorf("piv://%s: read certificate: %w", ref.Slot, err)
	}

	auth := piv.KeyAuth{}
	if ref.PIN != "" {
		auth.PIN = ref.PIN
	} else {
		auth.PIN = piv.DefaultPIN
	}

	priv, err := yk.PrivateKey(slot, cert.PublicKey, auth)
	if err != nil {
		yk.Close()
		return nil, "", fmt.Errorf("piv://%s: private key: %w", ref.Slot, err)
	}
	signer, ok := priv.(stdcrypto.Signer)
	if !ok {
		yk.Close()
		return nil, "", fmt.Errorf("piv://%s: key does not implement crypto.Signer", ref.Slot)
	}
	algo, err := algoFromSigner(signer)
	if err != nil {
		yk.Close()
		return nil, "", fmt.Errorf("piv://%s: %w", ref.Slot, err)
	}
	return &pivHolder{yk: yk, Signer: signer}, algo, nil
}

// pivHolder keeps the YubiKey handle alive for the lifetime of the signer.
type pivHolder struct {
	yk *piv.YubiKey
	stdcrypto.Signer
}

func (h *pivHolder) Close() error {
	if h.yk != nil {
		return h.yk.Close()
	}
	return nil
}

func pivSlot(name string) (piv.Slot, error) {
	switch strings.ToLower(name) {
	case "9a", "authentication":
		return piv.SlotAuthentication, nil
	case "9c", "signature", "digitalsignature":
		return piv.SlotSignature, nil
	case "9d", "keymanagement", "key-management":
		return piv.SlotKeyManagement, nil
	case "9e", "cardauthentication", "card-authentication":
		return piv.SlotCardAuthentication, nil
	default:
		return piv.Slot{}, fmt.Errorf("piv://%s: unknown slot (use 9a/9c/9d/9e)", name)
	}
}
