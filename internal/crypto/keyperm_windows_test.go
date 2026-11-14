//go:build windows

package crypto

import (
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

func TestSaveKeyFile_WindowsACL(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "secret.key")
	key, err := GenerateKey(AES256KeySize)
	if err != nil {
		t.Fatal(err)
	}
	if err := SaveKeyFile(keyPath, key); err != nil {
		t.Fatalf("SaveKeyFile: %v", err)
	}

	token := windows.GetCurrentProcessToken()
	tu, err := token.GetTokenUser()
	if err != nil {
		t.Fatalf("token user: %v", err)
	}
	sid := tu.User.Sid.String()

	sd, err := windows.GetNamedSecurityInfo(keyPath, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		t.Fatalf("GetNamedSecurityInfo: %v", err)
	}
	sddl := sd.String()
	if !strings.Contains(sddl, sid) {
		t.Fatalf("DACL missing current user SID %s: %s", sid, sddl)
	}
	// WD = Everyone, BU = Builtin\Users — must not appear after lock-down.
	for _, bad := range []string{";WD)", ";BU)", ";;WD)", ";;BU)"} {
		if strings.Contains(sddl, bad) {
			t.Fatalf("DACL grants broad group access (%s): %s", bad, sddl)
		}
	}
	if !strings.Contains(sddl, "D:P") && !strings.Contains(sddl, "D:PAI") && !strings.HasPrefix(sddl, "D:") {
		t.Fatalf("unexpected SDDL (want protected DACL): %s", sddl)
	}
}

func TestSaveKeyPEM_WindowsACL(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "priv.pem")
	if err := SaveKeyPEM(path, []byte("-----BEGIN PRIVATE KEY-----\nAA==\n-----END PRIVATE KEY-----\n"), 0o600); err != nil {
		t.Fatalf("SaveKeyPEM: %v", err)
	}
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		t.Fatalf("GetNamedSecurityInfo: %v", err)
	}
	sddl := sd.String()
	token := windows.GetCurrentProcessToken()
	tu, err := token.GetTokenUser()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sddl, tu.User.Sid.String()) {
		t.Fatalf("PEM key DACL missing current user: %s", sddl)
	}
}
