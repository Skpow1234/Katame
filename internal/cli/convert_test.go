package cli

import (
	"crypto/ed25519"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Skpow1234/Katame/internal/util"
)

func TestConvertCmd_LegacyEd25519(t *testing.T) {
	dir := t.TempDir()
	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	in := filepath.Join(dir, "legacy.key")
	out := filepath.Join(dir, "converted.pem")
	body := "ed25519-priv:" + util.B64Encode(priv) + "\n"
	if err := os.WriteFile(in, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}

	root := NewRootCmd()
	root.SetArgs([]string{"convert", "--in", in, "--out", out})
	if err := root.Execute(); err != nil {
		t.Fatalf("convert: %v", err)
	}
	pem, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pem), "BEGIN PRIVATE KEY") {
		t.Fatalf("expected PKCS#8 PEM, got: %s", pem)
	}

	root = NewRootCmd()
	root.SetArgs([]string{"convert", "--in", in, "--out", out})
	if err := root.Execute(); err == nil {
		t.Fatal("expected refuse overwrite without --force")
	}
	root = NewRootCmd()
	root.SetArgs([]string{"convert", "--in", in, "--out", out, "--force"})
	if err := root.Execute(); err != nil {
		t.Fatalf("convert --force: %v", err)
	}
}
