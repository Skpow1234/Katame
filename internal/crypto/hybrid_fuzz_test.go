package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func writeHybridKeyPair(t *testing.T, scheme string) (privPath, pubPath string) {
	t.Helper()
	privPEM, pubPEM, err := GenerateHybridKeys(scheme)
	if err != nil {
		t.Fatalf("GenerateHybridKeys(%s): %v", scheme, err)
	}
	dir := t.TempDir()
	privPath = filepath.Join(dir, "priv.pem")
	pubPath = filepath.Join(dir, "pub.pem")
	if err := SaveKeyPEM(privPath, privPEM, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := SaveKeyPEM(pubPath, pubPEM, 0o644); err != nil {
		t.Fatal(err)
	}
	return privPath, pubPath
}

func fuzzHybridDecapsulate(f *testing.F, scheme string) {
	privPEM, pubPEM, err := GenerateHybridKeys(scheme)
	if err != nil {
		f.Fatal(err)
	}
	dir, err := os.MkdirTemp("", "fuzz-hybrid-*")
	if err != nil {
		f.Fatal(err)
	}
	f.Cleanup(func() { _ = os.RemoveAll(dir) })
	privPath := filepath.Join(dir, "priv.pem")
	pubPath := filepath.Join(dir, "pub.pem")
	if err := SaveKeyPEM(privPath, privPEM, 0o600); err != nil {
		f.Fatal(err)
	}
	if err := SaveKeyPEM(pubPath, pubPEM, 0o644); err != nil {
		f.Fatal(err)
	}

	res, err := HybridEncapsulate(scheme, pubPath)
	if err != nil {
		f.Fatal(err)
	}
	f.Add(res.EphemeralPublicKey, res.WrappedDEK)
	f.Add([]byte{}, []byte{})
	f.Add(make([]byte, 8), make([]byte, 8))

	f.Fuzz(func(t *testing.T, eph, wrap []byte) {
		_, _ = HybridDecapsulate(scheme, privPath, eph, wrap)
	})
}

func fuzzHybridRoundTrip(f *testing.F, scheme string) {
	f.Add([]byte{1})
	f.Fuzz(func(t *testing.T, _ []byte) {
		privPath, pubPath := writeHybridKeyPair(t, scheme)
		res, err := HybridEncapsulate(scheme, pubPath)
		if err != nil {
			t.Fatalf("encapsulate: %v", err)
		}
		dek, err := HybridDecapsulate(scheme, privPath, res.EphemeralPublicKey, res.WrappedDEK)
		if err != nil {
			t.Fatalf("decapsulate: %v", err)
		}
		if !KeyEqual(dek, res.DEK) {
			t.Fatal("DEK mismatch")
		}
		Zeroize(dek)
		Zeroize(res.DEK)
	})
}

// FuzzHybridX25519Decapsulate (M35).
func FuzzHybridX25519Decapsulate(f *testing.F) {
	fuzzHybridDecapsulate(f, HybridX25519AES256GCM)
}

// FuzzHybridX25519RoundTrip (M35).
func FuzzHybridX25519RoundTrip(f *testing.F) {
	fuzzHybridRoundTrip(f, HybridX25519AES256GCM)
}

// FuzzHybridECIESP256Decapsulate (M36).
func FuzzHybridECIESP256Decapsulate(f *testing.F) {
	fuzzHybridDecapsulate(f, HybridECIESP256)
}

// FuzzHybridRSAOAEP2048Decapsulate (M37).
func FuzzHybridRSAOAEP2048Decapsulate(f *testing.F) {
	fuzzHybridDecapsulate(f, HybridRSAOAEP2048)
}

// FuzzHybridMLKEM768Decapsulate (M38).
func FuzzHybridMLKEM768Decapsulate(f *testing.F) {
	fuzzHybridDecapsulate(f, HybridMLKEM768)
}
