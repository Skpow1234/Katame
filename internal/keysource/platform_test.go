package keysource

import (
	"os"
	"path/filepath"
	"testing"

	vpcrypto "github.com/Skpow1234/Katame/internal/crypto"
	"github.com/Skpow1234/Katame/internal/util"
)

func TestPlatformSchemes_ReportsFour(t *testing.T) {
	schemes := PlatformSchemes()
	if len(schemes) != 4 {
		t.Fatalf("expected 4 platform schemes, got %d", len(schemes))
	}
}

func TestPKCS11_FileShim_Symmetric(t *testing.T) {
	dir := t.TempDir()
	key := make([]byte, vpcrypto.AES256KeySize)
	for i := range key {
		key[i] = byte(i + 1)
	}
	body := vpcrypto.KeyFilePrefix + util.B64Encode(key)
	if err := os.WriteFile(filepath.Join(dir, "demo.key"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("KATAME_PKCS11_DIR", dir)

	got, err := ResolveSymmetricKey("pkcs11://demo")
	if err != nil {
		t.Fatalf("ResolveSymmetricKey: %v", err)
	}
	if len(got) != vpcrypto.AES256KeySize {
		t.Fatalf("key len %d", len(got))
	}
	for i := range key {
		if got[i] != key[i] {
			t.Fatalf("mismatch at %d", i)
		}
	}
}

func TestPIV_FileShim_Private(t *testing.T) {
	dir := t.TempDir()
	pem := []byte("-----BEGIN PRIVATE KEY-----\nMII=\n-----END PRIVATE KEY-----\n")
	if err := os.WriteFile(filepath.Join(dir, "9a.pem"), pem, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("KATAME_PIV_DIR", dir)

	got, err := ResolvePrivateKey("piv://9a")
	if err != nil {
		t.Fatalf("ResolvePrivateKey: %v", err)
	}
	if string(got) != string(pem) {
		t.Fatalf("pem mismatch")
	}
}

func TestPKCS11_RequiresDir(t *testing.T) {
	t.Setenv("KATAME_PKCS11_DIR", "")
	_, err := ResolveSymmetricKey("pkcs11://missing")
	if err == nil {
		t.Fatal("expected error without KATAME_PKCS11_DIR")
	}
}
