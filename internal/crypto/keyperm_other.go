//go:build !windows

package crypto

// lockDownKeyFile is a no-op on non-Windows; Unix modes from WriteFile apply.
func lockDownKeyFile(path string) error { return nil }
