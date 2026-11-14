//go:build !windows

package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveKeyFile_UnixMode(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "secret.key")
	key, err := GenerateKey(AES256KeySize)
	if err != nil {
		t.Fatal(err)
	}
	if err := SaveKeyFile(keyPath, key); err != nil {
		t.Fatalf("SaveKeyFile: %v", err)
	}
	st, err := os.Stat(keyPath)
	if err != nil {
		t.Fatal(err)
	}
	if perm := st.Mode().Perm(); perm&0o077 != 0 {
		t.Fatalf("key file group/other bits set: %04o", perm)
	}
}
