package katame

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/Skpow1234/Katame/internal/bundle"
	"github.com/Skpow1234/Katame/internal/crypto"
)

// VerifyOptions controls a single Verify call.
//
// Exactly one of {PublicKey, PublicKeyPath} must be set. The bundle's
// signature is verified against the canonical manifest plus the SHA-256 of
// the payload.
type VerifyOptions struct {
	BundlePath    string
	PublicKey     []byte
	PublicKeyPath string
}

// VerifyResult is what Verify returns. Valid == true means the bundle's
// detached signature was successfully verified.
type VerifyResult struct {
	Valid     bool
	Algorithm string
	SignedAt  string
	Manifest  *Manifest
}

// Verify validates the detached signature of a .katame bundle.
func Verify(opts VerifyOptions) (*VerifyResult, error) {
	if opts.BundlePath == "" {
		return nil, errors.New("katame.Verify: BundlePath is required")
	}
	if opts.PublicKey == nil && opts.PublicKeyPath == "" {
		return nil, errors.New("katame.Verify: PublicKey or PublicKeyPath is required")
	}
	if opts.PublicKey != nil && opts.PublicKeyPath != "" {
		return nil, errors.New("katame.Verify: PublicKey and PublicKeyPath are mutually exclusive")
	}

	br, err := bundle.Read(opts.BundlePath)
	if err != nil {
		return nil, fmt.Errorf("katame.Verify: read: %w", err)
	}
	return verifyParsed(br, opts.PublicKey, opts.PublicKeyPath)
}

// VerifyBytes is the in-memory counterpart of Verify, accepting the raw
// bundle bytes directly. Browser / WASM friendly.
func VerifyBytes(bundleBytes []byte, opts VerifyOptions) (*VerifyResult, error) {
	if opts.PublicKey == nil && opts.PublicKeyPath == "" {
		return nil, errors.New("katame.VerifyBytes: PublicKey or PublicKeyPath is required")
	}
	br, err := bundle.ReadBytes(bundleBytes)
	if err != nil {
		return nil, fmt.Errorf("katame.VerifyBytes: read: %w", err)
	}
	return verifyParsed(br, opts.PublicKey, opts.PublicKeyPath)
}

func verifyParsed(br *bundle.ReadResult, pubBytes []byte, pubPath string) (*VerifyResult, error) {
	var err error
	if br.Signature == nil {
		return nil, errors.New("katame.Verify: bundle has no signature")
	}

	var (
		pub  any
		algo string
	)
	if pubPath != "" {
		pub, algo, err = crypto.LoadAnyPublicKey(pubPath)
	} else {
		pub, algo, err = crypto.ParseAnyPublicKey(pubBytes)
	}
	if err != nil {
		return nil, fmt.Errorf("katame.Verify: load public key: %w", err)
	}
	if br.Manifest.SignatureAlgo != nil && *br.Manifest.SignatureAlgo != "" {
		algo = *br.Manifest.SignatureAlgo
	}

	canonical, err := bundle.CanonicalManifest(br.Manifest)
	if err != nil {
		return nil, fmt.Errorf("katame.Verify: canonicalize: %w", err)
	}
	payloadHash, err := crypto.HashReader(bytes.NewReader(br.Ciphertext), "sha256")
	if err != nil {
		return nil, fmt.Errorf("katame.Verify: hash payload: %w", err)
	}
	sigMsg := crypto.BuildSigningMessage(canonical, payloadHash)

	valid, err := crypto.VerifySignature(pub, algo, sigMsg, br.Signature)
	if err != nil {
		return nil, fmt.Errorf("katame.Verify: verify: %w", err)
	}

	signedAt := ""
	if br.Manifest.SignedAt != nil {
		signedAt = *br.Manifest.SignedAt
	}
	return &VerifyResult{
		Valid:     valid,
		Algorithm: algo,
		SignedAt:  signedAt,
		Manifest:  br.Manifest,
	}, nil
}
