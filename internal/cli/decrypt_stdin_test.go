package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestDecryptCmd_Stdin(t *testing.T) {
	dir := t.TempDir()
	inFile := filepath.Join(dir, "secret.txt")
	bundleFile := filepath.Join(dir, "secret.katame")
	keyFile := filepath.Join(dir, "secret.key")
	outFile := filepath.Join(dir, "recovered.txt")
	original := []byte("stdin decrypt payload")
	if err := os.WriteFile(inFile, original, 0o600); err != nil {
		t.Fatal(err)
	}

	root := NewRootCmd()
	root.SetArgs([]string{"protect", "--in", inFile, "--out", bundleFile, "--key-out", keyFile})
	if err := root.Execute(); err != nil {
		t.Fatalf("protect: %v", err)
	}
	bundleBytes, err := os.ReadFile(bundleFile)
	if err != nil {
		t.Fatal(err)
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldIn := os.Stdin
	os.Stdin = r
	t.Cleanup(func() { os.Stdin = oldIn })
	go func() {
		_, _ = io.Copy(w, bytes.NewReader(bundleBytes))
		_ = w.Close()
	}()

	root = NewRootCmd()
	root.SetArgs([]string{"decrypt", "--stdin", "--out", outFile, "--key", keyFile})
	if err := root.Execute(); err != nil {
		t.Fatalf("decrypt --stdin: %v", err)
	}
	got, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(original) {
		t.Fatalf("got %q want %q", got, original)
	}
}
