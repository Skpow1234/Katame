//go:build !windows

package keysource

import "fmt"

const (
	dpapiAvailable = false
	dpapiDetails   = "Windows-only; rebuild on windows or use file:// / env://"
)

func resolveDPAPISymmetric(value string) ([]byte, error) {
	return nil, fmt.Errorf("dpapi:// is only available on Windows (got %q)", value)
}

func resolveDPAPIPrivate(value string) ([]byte, error) {
	return nil, fmt.Errorf("dpapi:// is only available on Windows (got %q)", value)
}
