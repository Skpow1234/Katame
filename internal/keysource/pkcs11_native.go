//go:build cgo && katame_pkcs11

package keysource

import (
	stdcrypto "crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/asn1"
	"fmt"
	"io"
	"math/big"

	"github.com/miekg/pkcs11"
)

const nativePKCS11Enabled = true

func resolveNativePKCS11Symmetric(ref pkcs11Ref) ([]byte, error) {
	mod := pkcs11ModulePath()
	if mod == "" {
		return nil, fmt.Errorf("pkcs11://%s: set KATAME_PKCS11_MODULE to SoftHSM libsofthsm2.so", ref.Label)
	}
	ctx, session, err := openPKCS11Session(mod, ref)
	if err != nil {
		return nil, err
	}
	defer closePKCS11Session(ctx, session)

	obj, err := findPKCS11Object(ctx, session, []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, ref.Label),
	})
	if err != nil {
		return nil, fmt.Errorf("pkcs11://%s: %w", ref.Label, err)
	}
	attrs, err := ctx.GetAttributeValue(session, obj, []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_VALUE, nil),
	})
	if err != nil {
		return nil, fmt.Errorf("pkcs11://%s: read key value (is the key extractable?): %w", ref.Label, err)
	}
	if len(attrs) == 0 || len(attrs[0].Value) == 0 {
		return nil, fmt.Errorf("pkcs11://%s: empty key value", ref.Label)
	}
	key := attrs[0].Value
	if len(key) != 32 {
		return nil, fmt.Errorf("pkcs11://%s: want 32-byte AES key, got %d", ref.Label, len(key))
	}
	return key, nil
}

func resolveNativePKCS11Signer(ref pkcs11Ref) (stdcrypto.Signer, string, error) {
	mod := pkcs11ModulePath()
	if mod == "" {
		return nil, "", fmt.Errorf("pkcs11://%s: set KATAME_PKCS11_MODULE to SoftHSM libsofthsm2.so", ref.Label)
	}
	ctx, session, err := openPKCS11Session(mod, ref)
	if err != nil {
		return nil, "", err
	}

	obj, err := findPKCS11Object(ctx, session, []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, ref.Label),
	})
	if err != nil {
		closePKCS11Session(ctx, session)
		return nil, "", fmt.Errorf("pkcs11://%s: %w", ref.Label, err)
	}

	pub, algo, err := loadPKCS11Public(ctx, session, ref.Label)
	if err != nil {
		closePKCS11Session(ctx, session)
		return nil, "", err
	}

	return &pkcs11Signer{
		ctx:     ctx,
		session: session,
		obj:     obj,
		pub:     pub,
		algo:    algo,
	}, algo, nil
}

type pkcs11Signer struct {
	ctx     *pkcs11.Ctx
	session pkcs11.SessionHandle
	obj     pkcs11.ObjectHandle
	pub     stdcrypto.PublicKey
	algo    string
}

func (s *pkcs11Signer) Public() stdcrypto.PublicKey { return s.pub }

func (s *pkcs11Signer) Sign(rand io.Reader, digest []byte, opts stdcrypto.SignerOpts) ([]byte, error) {
	_ = rand
	switch s.algo {
	case "ecdsa-p256", "ecdsa-p384":
		mech := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA, nil)}
		if err := s.ctx.SignInit(s.session, mech, s.obj); err != nil {
			return nil, fmt.Errorf("pkcs11 SignInit ECDSA: %w", err)
		}
		raw, err := s.ctx.Sign(s.session, digest)
		if err != nil {
			return nil, fmt.Errorf("pkcs11 Sign ECDSA: %w", err)
		}
		return ecdsaRawToASN1(raw)
	case "rsa-pss-2048", "rsa-pss-4096":
		hashLen := 32
		if opts != nil {
			if h := opts.HashFunc(); h != 0 {
				hashLen = h.Size()
			}
		}
		params := pkcs11.NewPSSParams(pkcs11.CKM_SHA256, pkcs11.CKG_MGF1_SHA256, uint(hashLen))
		mech := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_PSS, params)}
		if err := s.ctx.SignInit(s.session, mech, s.obj); err != nil {
			return nil, fmt.Errorf("pkcs11 SignInit RSA-PSS: %w", err)
		}
		return s.ctx.Sign(s.session, digest)
	default:
		return nil, fmt.Errorf("pkcs11 signer: unsupported algo %q", s.algo)
	}
}

func ecdsaRawToASN1(raw []byte) ([]byte, error) {
	if len(raw)%2 != 0 || len(raw) == 0 {
		return nil, fmt.Errorf("invalid ECDSA signature length %d", len(raw))
	}
	half := len(raw) / 2
	r := new(big.Int).SetBytes(raw[:half])
	s := new(big.Int).SetBytes(raw[half:])
	return asn1.Marshal(struct{ R, S *big.Int }{r, s})
}

func loadPKCS11Public(ctx *pkcs11.Ctx, session pkcs11.SessionHandle, label string) (stdcrypto.PublicKey, string, error) {
	obj, err := findPKCS11Object(ctx, session, []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
	})
	if err != nil {
		return nil, "", fmt.Errorf("pkcs11://%s: public key: %w", label, err)
	}
	typeAttrs, err := ctx.GetAttributeValue(session, obj, []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, nil),
	})
	if err != nil || len(typeAttrs) == 0 {
		return nil, "", fmt.Errorf("pkcs11://%s: read key type: %w", label, err)
	}
	keyType := bytesToUlong(typeAttrs[0].Value)
	switch keyType {
	case pkcs11.CKK_EC:
		attrs, err := ctx.GetAttributeValue(session, obj, []*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_EC_POINT, nil),
			pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, nil),
		})
		if err != nil || len(attrs) < 2 {
			return nil, "", fmt.Errorf("pkcs11://%s: read EC public: %w", label, err)
		}
		curve, algo, err := curveFromECParams(attrs[1].Value)
		if err != nil {
			return nil, "", err
		}
		x, y, err := parseECPoint(attrs[0].Value, curve)
		if err != nil {
			return nil, "", err
		}
		return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, algo, nil
	case pkcs11.CKK_RSA:
		attrs, err := ctx.GetAttributeValue(session, obj, []*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_MODULUS, nil),
			pkcs11.NewAttribute(pkcs11.CKA_PUBLIC_EXPONENT, nil),
		})
		if err != nil || len(attrs) < 2 {
			return nil, "", fmt.Errorf("pkcs11://%s: read RSA public: %w", label, err)
		}
		n := new(big.Int).SetBytes(attrs[0].Value)
		e := int(new(big.Int).SetBytes(attrs[1].Value).Int64())
		pub := &rsa.PublicKey{N: n, E: e}
		algo := "rsa-pss-4096"
		if n.BitLen() <= 2048 {
			algo = "rsa-pss-2048"
		}
		return pub, algo, nil
	default:
		return nil, "", fmt.Errorf("pkcs11://%s: unsupported key type %d", label, keyType)
	}
}

func curveFromECParams(params []byte) (elliptic.Curve, string, error) {
	// Named curve OIDs: secp256r1 / secp384r1
	var oid asn1.ObjectIdentifier
	if _, err := asn1.Unmarshal(params, &oid); err != nil {
		return nil, "", fmt.Errorf("parse EC params: %w", err)
	}
	switch {
	case oid.Equal(asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}):
		return elliptic.P256(), "ecdsa-p256", nil
	case oid.Equal(asn1.ObjectIdentifier{1, 3, 132, 0, 34}):
		return elliptic.P384(), "ecdsa-p384", nil
	default:
		return nil, "", fmt.Errorf("unsupported EC curve OID %v", oid)
	}
}

func parseECPoint(ecPoint []byte, curve elliptic.Curve) (*big.Int, *big.Int, error) {
	// SoftHSM stores EC_POINT as OCTET STRING wrapping uncompressed point.
	var point []byte
	if _, err := asn1.Unmarshal(ecPoint, &point); err != nil {
		point = ecPoint
	}
	x, y := elliptic.Unmarshal(curve, point)
	if x == nil {
		return nil, nil, fmt.Errorf("invalid EC point")
	}
	return x, y, nil
}

func bytesToUlong(b []byte) uint {
	var n uint
	for _, v := range b {
		n = (n << 8) | uint(v)
	}
	return n
}

func openPKCS11Session(module string, ref pkcs11Ref) (*pkcs11.Ctx, pkcs11.SessionHandle, error) {
	ctx := pkcs11.New(module)
	if ctx == nil {
		return nil, 0, fmt.Errorf("pkcs11: failed to load module %q", module)
	}
	if err := ctx.Initialize(); err != nil {
		ctx.Destroy()
		return nil, 0, fmt.Errorf("pkcs11 Initialize: %w", err)
	}
	slots, err := ctx.GetSlotList(true)
	if err != nil || len(slots) == 0 {
		ctx.Finalize()
		ctx.Destroy()
		return nil, 0, fmt.Errorf("pkcs11: no token slots: %w", err)
	}
	slot := slots[0]
	if ref.Slot >= 0 {
		found := false
		for _, s := range slots {
			if int(s) == ref.Slot {
				slot = s
				found = true
				break
			}
		}
		if !found {
			ctx.Finalize()
			ctx.Destroy()
			return nil, 0, fmt.Errorf("pkcs11: slot %d not found", ref.Slot)
		}
	}
	session, err := ctx.OpenSession(slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		ctx.Finalize()
		ctx.Destroy()
		return nil, 0, fmt.Errorf("pkcs11 OpenSession: %w", err)
	}
	pin := ref.PIN
	if pin == "" {
		pin = "1234" // SoftHSM default user PIN in many tutorials; override via ?pin= or KATAME_PKCS11_PIN
	}
	if err := ctx.Login(session, pkcs11.CKU_USER, pin); err != nil {
		_ = ctx.CloseSession(session)
		ctx.Finalize()
		ctx.Destroy()
		return nil, 0, fmt.Errorf("pkcs11 Login: %w", err)
	}
	return ctx, session, nil
}

func closePKCS11Session(ctx *pkcs11.Ctx, session pkcs11.SessionHandle) {
	_ = ctx.Logout(session)
	_ = ctx.CloseSession(session)
	_ = ctx.Finalize()
	ctx.Destroy()
}

func findPKCS11Object(ctx *pkcs11.Ctx, session pkcs11.SessionHandle, template []*pkcs11.Attribute) (pkcs11.ObjectHandle, error) {
	if err := ctx.FindObjectsInit(session, template); err != nil {
		return 0, err
	}
	objs, _, err := ctx.FindObjects(session, 1)
	if err != nil {
		_ = ctx.FindObjectsFinal(session)
		return 0, err
	}
	if err := ctx.FindObjectsFinal(session); err != nil {
		return 0, err
	}
	if len(objs) == 0 {
		return 0, fmt.Errorf("object not found")
	}
	return objs[0], nil
}
