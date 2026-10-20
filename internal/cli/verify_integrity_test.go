package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyIntegrity_KeyFile(t *testing.T) {
	dir := t.TempDir()
	inFile := filepath.Join(dir, "data.txt")
	outKatame := filepath.Join(dir, "data.katame")
	keyFile := filepath.Join(dir, "data.key")

	payload := []byte("integrity check data")
	os.WriteFile(inFile, payload, 0o600)

	// Protect.
	root := NewRootCmd()
	root.SetArgs([]string{
		"protect",
		"--in", inFile,
		"--out", outKatame,
		"--key-out", keyFile,
	})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}

	// Verify integrity.
	root2 := NewRootCmd()
	root2.SetArgs([]string{
		"verify-integrity",
		"--in", outKatame,
		"--key", keyFile,
	})
	if err := root2.Execute(); err != nil {
		t.Fatalf("verify-integrity: %v", err)
	}
}

func TestVerifyIntegrity_Password(t *testing.T) {
	dir := t.TempDir()
	inFile := filepath.Join(dir, "data.txt")
	outKatame := filepath.Join(dir, "data.katame")

	os.WriteFile(inFile, []byte("password integrity"), 0o600)

	root := NewRootCmd()
	root.SetArgs([]string{
		"protect",
		"--in", inFile,
		"--out", outKatame,
		"--password", "testpass123",
		"--kdf", "pbkdf2-sha256",
	})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}

	root2 := NewRootCmd()
	root2.SetArgs([]string{
		"verify-integrity",
		"--in", outKatame,
		"--password", "testpass123",
	})
	if err := root2.Execute(); err != nil {
		t.Fatalf("verify-integrity with password: %v", err)
	}
}

func TestVerifyIntegrity_WithCompression(t *testing.T) {
	dir := t.TempDir()
	inFile := filepath.Join(dir, "data.txt")
	outKatame := filepath.Join(dir, "data.katame")
	keyFile := filepath.Join(dir, "data.key")

	os.WriteFile(inFile, []byte("compressed integrity check"), 0o600)

	root := NewRootCmd()
	root.SetArgs([]string{
		"protect",
		"--in", inFile,
		"--out", outKatame,
		"--key-out", keyFile,
		"--compress", "gzip",
	})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}

	root2 := NewRootCmd()
	root2.SetArgs([]string{
		"verify-integrity",
		"--in", outKatame,
		"--key", keyFile,
	})
	if err := root2.Execute(); err != nil {
		t.Fatalf("verify-integrity with compression: %v", err)
	}
}

func TestVerifyIntegrity_MissingKey(t *testing.T) {
	root := NewRootCmd()
	root.SetArgs([]string{
		"verify-integrity",
		"--in", "test.katame",
	})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when no key method provided")
	}
}
