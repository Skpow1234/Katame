package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Skpow1234/Katame/internal/bundle"
	"github.com/Skpow1234/Katame/internal/crypto"
)

func TestSignTimestamp(t *testing.T) {
	dir := t.TempDir()

	// Generate a signing key pair.
	privPEM, pubPEM, err := crypto.GenerateSigningKeys(crypto.SignAlgoEd25519)
	if err != nil {
		t.Fatal(err)
	}
	privFile := filepath.Join(dir, "sign.key")
	pubFile := filepath.Join(dir, "sign.pub")
	crypto.SaveKeyPEM(privFile, privPEM, 0o600)
	crypto.SaveKeyPEM(pubFile, pubPEM, 0o644)

	inFile := filepath.Join(dir, "data.txt")
	outKatame := filepath.Join(dir, "data.katame")
	keyFile := filepath.Join(dir, "data.key")

	os.WriteFile(inFile, []byte("timestamp test"), 0o600)

	// Protect with signing.
	root := NewRootCmd()
	root.SetArgs([]string{
		"protect",
		"--in", inFile,
		"--out", outKatame,
		"--key-out", keyFile,
		"--sign",
		"--signing-priv", privFile,
	})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}

	// Read manifest and check timestamp.
	br, err := bundle.Read(outKatame)
	if err != nil {
		t.Fatal(err)
	}
	if br.Manifest.SignedAt == nil {
		t.Fatal("expected signed_at timestamp")
	}
	if *br.Manifest.SignedAt == "" {
		t.Error("signed_at is empty")
	}

	// Verify still works.
	root2 := NewRootCmd()
	root2.SetArgs([]string{
		"verify",
		"--in", outKatame,
		"--pubkey", pubFile,
	})
	if err := root2.Execute(); err != nil {
		t.Fatalf("verify: %v", err)
	}
}

func TestSignCmd_Timestamp(t *testing.T) {
	dir := t.TempDir()

	privPEM, pubPEM, _ := crypto.GenerateSigningKeys(crypto.SignAlgoEd25519)
	privFile := filepath.Join(dir, "sign.key")
	pubFile := filepath.Join(dir, "sign.pub")
	crypto.SaveKeyPEM(privFile, privPEM, 0o600)
	crypto.SaveKeyPEM(pubFile, pubPEM, 0o644)

	inFile := filepath.Join(dir, "data.txt")
	outKatame := filepath.Join(dir, "data.katame")
	keyFile := filepath.Join(dir, "data.key")

	os.WriteFile(inFile, []byte("sign command timestamp"), 0o600)

	// Protect without signing.
	root := NewRootCmd()
	root.SetArgs([]string{
		"protect",
		"--in", inFile,
		"--out", outKatame,
		"--key-out", keyFile,
	})
	root.Execute()

	// Sign separately.
	root2 := NewRootCmd()
	root2.SetArgs([]string{
		"sign",
		"--in", outKatame,
		"--signing-priv", privFile,
	})
	if err := root2.Execute(); err != nil {
		t.Fatalf("sign: %v", err)
	}

	// Read manifest and check timestamp.
	br, err := bundle.Read(outKatame)
	if err != nil {
		t.Fatal(err)
	}
	if br.Manifest.SignedAt == nil {
		t.Fatal("expected signed_at from sign command")
	}

	// Verify.
	root3 := NewRootCmd()
	root3.SetArgs([]string{
		"verify",
		"--in", outKatame,
		"--pubkey", pubFile,
	})
	if err := root3.Execute(); err != nil {
		t.Fatalf("verify after sign: %v", err)
	}
}
