package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRefuseIfExists(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "exists.txt")
	if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := refuseIfExists(p, false); err == nil {
		t.Fatal("expected refuse without --force")
	} else if !strings.Contains(err.Error(), "--force") {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := refuseIfExists(p, true); err != nil {
		t.Fatalf("force should allow overwrite: %v", err)
	}
	missing := filepath.Join(dir, "missing.txt")
	if err := refuseIfExists(missing, false); err != nil {
		t.Fatalf("missing path should be ok: %v", err)
	}
	if err := refuseIfExists("", false); err != nil {
		t.Fatalf("empty path should be ok: %v", err)
	}
}

func TestProtectCmd_ForceOverwrite(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.txt")
	out := filepath.Join(dir, "out.katame")
	key := filepath.Join(dir, "out.key")
	if err := os.WriteFile(in, []byte("force-me"), 0o600); err != nil {
		t.Fatal(err)
	}
	root := NewRootCmd()
	root.SetArgs([]string{"protect", "--in", in, "--out", out, "--key-out", key})
	if err := root.Execute(); err != nil {
		t.Fatalf("first protect: %v", err)
	}
	root = NewRootCmd()
	root.SetArgs([]string{"protect", "--in", in, "--out", out, "--key-out", key})
	if err := root.Execute(); err == nil {
		t.Fatal("expected refuse overwrite without --force")
	}
	root = NewRootCmd()
	root.SetArgs([]string{"protect", "--in", in, "--out", out, "--key-out", key, "--force"})
	if err := root.Execute(); err != nil {
		t.Fatalf("protect --force: %v", err)
	}
}
