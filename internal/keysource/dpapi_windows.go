//go:build windows

package keysource

import (
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/Skpow1234/Katame/internal/util"
	"golang.org/x/sys/windows"
)

const (
	dpapiAvailable = true
	dpapiDetails   = "CryptProtectData/CryptUnprotectData (current-user scope)"
)

// dpapi://path/to/blob  — file contains DPAPI-protected ciphertext of a Katame key.
// dpapi://b64:<base64>  — inline DPAPI blob.
func resolveDPAPISymmetric(value string) ([]byte, error) {
	blob, err := loadDPAPIBlob(value)
	if err != nil {
		return nil, err
	}
	plain, err := dpapiUnprotect(blob)
	if err != nil {
		return nil, err
	}
	return parseSymmetricKeyBytes(plain)
}

func resolveDPAPIPrivate(value string) ([]byte, error) {
	blob, err := loadDPAPIBlob(value)
	if err != nil {
		return nil, err
	}
	return dpapiUnprotect(blob)
}

// ProtectDPAPI encrypts plaintext with DPAPI (current user) for storage under a dpapi:// URI.
func ProtectDPAPI(plaintext []byte) ([]byte, error) {
	return dpapiProtect(plaintext)
}

func loadDPAPIBlob(value string) ([]byte, error) {
	if strings.HasPrefix(value, "b64:") {
		return util.B64Decode(strings.TrimPrefix(value, "b64:"))
	}
	return os.ReadFile(fileValue(value))
}

func dpapiProtect(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("dpapi protect: empty plaintext")
	}
	in := windows.DataBlob{Size: uint32(len(data)), Data: &data[0]}
	var out windows.DataBlob
	const cryptprotectUIForbidden = 0x1
	err := windows.CryptProtectData(&in, nil, nil, 0, nil, cryptprotectUIForbidden, &out)
	if err != nil {
		return nil, fmt.Errorf("dpapi protect: %w", err)
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	return blobBytes(out), nil
}

func dpapiUnprotect(blob []byte) ([]byte, error) {
	if len(blob) == 0 {
		return nil, fmt.Errorf("dpapi: empty blob")
	}
	in := windows.DataBlob{Size: uint32(len(blob)), Data: &blob[0]}
	var out windows.DataBlob
	const cryptprotectUIForbidden = 0x1
	err := windows.CryptUnprotectData(&in, nil, nil, 0, nil, cryptprotectUIForbidden, &out)
	if err != nil {
		return nil, fmt.Errorf("dpapi unprotect: %w", err)
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	return blobBytes(out), nil
}

func blobBytes(b windows.DataBlob) []byte {
	if b.Size == 0 || b.Data == nil {
		return nil
	}
	out := make([]byte, b.Size)
	copy(out, unsafe.Slice(b.Data, b.Size))
	return out
}
