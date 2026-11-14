//go:build windows

package crypto

import (
	"fmt"

	"golang.org/x/sys/windows"
)

// lockDownKeyFile replaces the file DACL with a protected ACL that grants
// full access only to the current process user (M108).
func lockDownKeyFile(path string) error {
	token := windows.GetCurrentProcessToken()
	tu, err := token.GetTokenUser()
	if err != nil {
		return fmt.Errorf("token user: %w", err)
	}
	sddl := "D:P(A;;FA;;;" + tu.User.Sid.String() + ")"
	sd, err := windows.SecurityDescriptorFromString(sddl)
	if err != nil {
		return fmt.Errorf("build key ACL: %w", err)
	}
	dacl, _, err := sd.DACL()
	if err != nil {
		return fmt.Errorf("key DACL: %w", err)
	}
	if err := windows.SetNamedSecurityInfo(
		path,
		windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION,
		nil, nil, dacl, nil,
	); err != nil {
		return fmt.Errorf("set key ACL: %w", err)
	}
	return nil
}
