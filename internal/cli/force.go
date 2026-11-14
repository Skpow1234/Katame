package cli

import (
	"fmt"
	"os"
)

// refuseIfExists returns an error when path already exists and force is false.
// Missing paths and empty path are allowed.
func refuseIfExists(path string, force bool) error {
	if path == "" {
		return nil
	}
	_, err := os.Stat(path)
	if err == nil {
		if force {
			return nil
		}
		return fmt.Errorf("refusing to overwrite existing path %q (use --force)", path)
	}
	if os.IsNotExist(err) {
		return nil
	}
	return fmt.Errorf("stat %q: %w", path, err)
}
