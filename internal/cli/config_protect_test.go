package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Skpow1234/Katame/internal/bundle"
	"github.com/Skpow1234/Katame/internal/config"
)

func TestConfigLoad_MalformedSurfaced(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(bad, []byte("cipher: [\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--config", bad, "doctor"})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "config load failed") {
		t.Fatalf("expected config load failure, got %v", err)
	}
}

func TestProtect_DryRun(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "secret.txt")
	if err := writeFileHelper(in, []byte("hello")); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "secret.katame")
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"protect", "--in", in, "--out", out, "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("protect dry-run: %v", err)
	}
	if _, err := os.Stat(out); err == nil {
		t.Fatal("dry-run should not write bundle")
	}
}

func TestProtect_ConfigCipherAndChunkSize(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "katame.yaml")
	if err := os.WriteFile(cfgPath, []byte("cipher: chacha20-poly1305\nchunk_size: 8192\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	in := filepath.Join(dir, "secret.txt")
	out := filepath.Join(dir, "secret.katame")
	if err := writeFileHelper(in, []byte("config-cipher")); err != nil {
		t.Fatal(err)
	}
	config.SetLoaded(nil)
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--config", cfgPath, "protect", "--in", in, "--out", out})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("protect with config: %v", err)
	}
	br, err := bundle.Read(out)
	if err != nil {
		t.Fatalf("read bundle: %v", err)
	}
	if br.Manifest.Encryption.AEAD != "chacha20-poly1305" {
		t.Fatalf("cipher=%q", br.Manifest.Encryption.AEAD)
	}
	if br.Manifest.Encryption.ChunkSize == nil || *br.Manifest.Encryption.ChunkSize != 8192 {
		t.Fatalf("chunk_size=%v", br.Manifest.Encryption.ChunkSize)
	}
}

func TestKeygen_ListAlgos(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"keygen", "--list-algos"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("keygen --list-algos: %v", err)
	}
}

func TestRewrap_DryRun(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{
		"rewrap", "--in", "missing.katame", "--kms-provider", "mock",
		"--to-kms-key-id", "k2", "--dry-run",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("rewrap dry-run: %v", err)
	}
}

func TestRotateKey_DryRun(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"rotate-key", "--in", "missing.katame", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("rotate-key dry-run: %v", err)
	}
}
