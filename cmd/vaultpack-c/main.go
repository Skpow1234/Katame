//go:build cgo && cshared

// Package main builds a C-shared library (libvaultpack.so / .dylib / .dll)
// exposing the VaultPack SDK to native callers (Python, Ruby, Java, etc).
//
// Build:
//
//	go build -buildmode=c-shared -tags cshared -o libvaultpack.so ./cmd/vaultpack-c
//
// The exported C ABI is intentionally small and JSON-based: every entry
// point takes a single null-terminated JSON string and returns a single
// malloc-allocated JSON string. The caller is responsible for releasing the
// returned string with vp_free.
package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"encoding/json"
	"unsafe"

	_ "github.com/Skpow1234/Vaultpack/internal/kms" // register aws/gcp/azure/mock providers
	"github.com/Skpow1234/Vaultpack/internal/util"
	"github.com/Skpow1234/Vaultpack/pkg/vaultpack"
)

func jsonReply(v any) *C.char {
	data, err := json.Marshal(v)
	if err != nil {
		data = []byte(`{"ok":false,"error":"jsonReply: marshal failed"}`)
	}
	return C.CString(string(data))
}

func errReply(msg string) *C.char {
	return jsonReply(map[string]any{"ok": false, "error": msg})
}

//export vp_version
func vp_version() *C.char {
	return C.CString(vaultpack.Version)
}

//export vp_free
func vp_free(p *C.char) {
	if p != nil {
		C.free(unsafe.Pointer(p))
	}
}

type recipientIn struct {
	PublicKeyPEM  string `json:"public_key_pem,omitempty"`
	PublicKeyPath string `json:"public_key_path,omitempty"`
}

//export vp_protect
func vp_protect(cOpts *C.char) *C.char {
	var in struct {
		InputPath      string        `json:"input_path,omitempty"`
		PlaintextB64   string        `json:"plaintext_b64,omitempty"`
		OutputPath     string        `json:"output_path,omitempty"`
		KeyB64         string        `json:"key_b64,omitempty"`
		Password       string        `json:"password,omitempty"`
		KDFAlgo        string        `json:"kdf_algo,omitempty"`
		Cipher         string        `json:"cipher,omitempty"`
		AADB64         string        `json:"aad_b64,omitempty"`
		InputName      string        `json:"input_name,omitempty"`
		Compress       string        `json:"compress,omitempty"`
		KMSProvider    string        `json:"kms_provider,omitempty"`
		KMSKeyID       string        `json:"kms_key_id,omitempty"`
		Recipients     []recipientIn `json:"recipients,omitempty"`
		SplitShares    int           `json:"split_shares,omitempty"`
		SplitThreshold int           `json:"split_threshold,omitempty"`
		HashAlgo       string        `json:"hash_algo,omitempty"`
	}
	if err := json.Unmarshal([]byte(C.GoString(cOpts)), &in); err != nil {
		return errReply("invalid JSON: " + err.Error())
	}

	opts := vaultpack.ProtectOptions{
		InputPath:      in.InputPath,
		OutputPath:     in.OutputPath,
		Password:       in.Password,
		KDFAlgo:        in.KDFAlgo,
		Cipher:         in.Cipher,
		InputName:      in.InputName,
		Compress:       in.Compress,
		KMSProvider:    in.KMSProvider,
		KMSKeyID:       in.KMSKeyID,
		SplitShares:    in.SplitShares,
		SplitThreshold: in.SplitThreshold,
		HashAlgo:       in.HashAlgo,
	}
	if in.PlaintextB64 != "" {
		pt, err := util.B64Decode(in.PlaintextB64)
		if err != nil {
			return errReply("decode plaintext_b64: " + err.Error())
		}
		opts.Plaintext = pt
	}
	if in.KeyB64 != "" {
		k, err := util.B64Decode(in.KeyB64)
		if err != nil {
			return errReply("decode key_b64: " + err.Error())
		}
		opts.Key = k
	}
	if in.AADB64 != "" {
		a, err := util.B64Decode(in.AADB64)
		if err != nil {
			return errReply("decode aad_b64: " + err.Error())
		}
		opts.AAD = a
	}
	for _, r := range in.Recipients {
		rec := vaultpack.Recipient{PublicKeyPath: r.PublicKeyPath}
		if r.PublicKeyPEM != "" {
			rec.PublicKeyPEM = []byte(r.PublicKeyPEM)
		}
		opts.Recipients = append(opts.Recipients, rec)
	}

	res, err := vaultpack.Protect(opts)
	if err != nil {
		return errReply(err.Error())
	}
	out := map[string]any{
		"ok":          true,
		"bundle_path": res.BundlePath,
		"manifest":    res.Manifest,
	}
	if res.GeneratedKey != nil {
		out["generated_key_b64"] = util.B64Encode(res.GeneratedKey)
	}
	if len(res.Shares) > 0 {
		shares := make([]map[string]any, len(res.Shares))
		for i, sh := range res.Shares {
			shares[i] = map[string]any{"index": sh.Index, "data_b64": util.B64Encode(sh.Data)}
		}
		out["shares"] = shares
	}
	return jsonReply(out)
}

//export vp_decrypt
func vp_decrypt(cOpts *C.char) *C.char {
	var in struct {
		InputPath      string `json:"input_path,omitempty"`
		BundleB64      string `json:"bundle_b64,omitempty"`
		OutputPath     string `json:"output_path,omitempty"`
		KeyB64         string `json:"key_b64,omitempty"`
		Password       string `json:"password,omitempty"`
		AADB64         string `json:"aad_b64,omitempty"`
		KMSProvider    string `json:"kms_provider,omitempty"`
		PrivateKeyPEM  string `json:"private_key_pem,omitempty"`
		PrivateKeyPath string `json:"private_key_path,omitempty"`
	}
	if err := json.Unmarshal([]byte(C.GoString(cOpts)), &in); err != nil {
		return errReply("invalid JSON: " + err.Error())
	}

	opts := vaultpack.DecryptOptions{
		InputPath:      in.InputPath,
		OutputPath:     in.OutputPath,
		Password:       in.Password,
		KMSProvider:    in.KMSProvider,
		PrivateKeyPath: in.PrivateKeyPath,
	}
	if in.BundleB64 != "" {
		b, err := util.B64Decode(in.BundleB64)
		if err != nil {
			return errReply("decode bundle_b64: " + err.Error())
		}
		opts.InputBytes = b
	}
	if in.KeyB64 != "" {
		k, err := util.B64Decode(in.KeyB64)
		if err != nil {
			return errReply("decode key_b64: " + err.Error())
		}
		opts.Key = k
	}
	if in.PrivateKeyPEM != "" {
		opts.PrivateKey = []byte(in.PrivateKeyPEM)
	}
	if in.AADB64 != "" {
		a, err := util.B64Decode(in.AADB64)
		if err != nil {
			return errReply("decode aad_b64: " + err.Error())
		}
		opts.AAD = a
	}

	res, err := vaultpack.Decrypt(opts)
	if err != nil {
		return errReply(err.Error())
	}
	out := map[string]any{
		"ok":       true,
		"manifest": res.Manifest,
	}
	if res.Plaintext != nil {
		out["plaintext_b64"] = util.B64Encode(res.Plaintext)
	}
	return jsonReply(out)
}

//export vp_inspect
func vp_inspect(cOpts *C.char) *C.char {
	raw := C.GoString(cOpts)
	// Accept either a raw path string or a JSON object for consistency.
	var in struct {
		InputPath string `json:"input_path,omitempty"`
		BundleB64 string `json:"bundle_b64,omitempty"`
		Path      string `json:"path,omitempty"`
	}
	path := raw
	if err := json.Unmarshal([]byte(raw), &in); err == nil && (in.InputPath != "" || in.BundleB64 != "" || in.Path != "") {
		if in.BundleB64 != "" {
			b, err := util.B64Decode(in.BundleB64)
			if err != nil {
				return errReply("decode bundle_b64: " + err.Error())
			}
			m, err := vaultpack.InspectBytes(b)
			if err != nil {
				return errReply(err.Error())
			}
			return jsonReply(map[string]any{"ok": true, "manifest": m})
		}
		path = in.InputPath
		if path == "" {
			path = in.Path
		}
	}
	m, err := vaultpack.Inspect(path)
	if err != nil {
		return errReply(err.Error())
	}
	return jsonReply(map[string]any{"ok": true, "manifest": m})
}

//export vp_sign
func vp_sign(cOpts *C.char) *C.char {
	var in struct {
		BundlePath     string `json:"bundle_path"`
		PrivateKeyPEM  string `json:"private_key_pem,omitempty"`
		PrivateKeyPath string `json:"private_key_path,omitempty"`
		Algo           string `json:"algo,omitempty"`
	}
	if err := json.Unmarshal([]byte(C.GoString(cOpts)), &in); err != nil {
		return errReply("invalid JSON: " + err.Error())
	}
	opts := vaultpack.SignOptions{
		BundlePath:     in.BundlePath,
		PrivateKeyPath: in.PrivateKeyPath,
		Algo:           in.Algo,
	}
	if in.PrivateKeyPEM != "" {
		opts.PrivateKey = []byte(in.PrivateKeyPEM)
	}
	res, err := vaultpack.SignBundle(opts)
	if err != nil {
		return errReply(err.Error())
	}
	return jsonReply(map[string]any{
		"ok":            true,
		"algorithm":     res.Algorithm,
		"signed_at":     res.SignedAt,
		"signature_b64": util.B64Encode(res.Signature),
	})
}

//export vp_verify
func vp_verify(cOpts *C.char) *C.char {
	var in struct {
		BundlePath    string `json:"bundle_path"`
		PublicKeyPEM  string `json:"public_key_pem,omitempty"`
		PublicKeyPath string `json:"public_key_path,omitempty"`
	}
	if err := json.Unmarshal([]byte(C.GoString(cOpts)), &in); err != nil {
		return errReply("invalid JSON: " + err.Error())
	}
	opts := vaultpack.VerifyOptions{
		BundlePath:    in.BundlePath,
		PublicKeyPath: in.PublicKeyPath,
	}
	if in.PublicKeyPEM != "" {
		opts.PublicKey = []byte(in.PublicKeyPEM)
	}
	res, err := vaultpack.Verify(opts)
	if err != nil {
		return errReply(err.Error())
	}
	return jsonReply(map[string]any{
		"ok":        true,
		"valid":     res.Valid,
		"algorithm": res.Algorithm,
		"signed_at": res.SignedAt,
		"manifest":  res.Manifest,
	})
}

//export vp_rewrap
func vp_rewrap(cOpts *C.char) *C.char {
	var in struct {
		InputPath   string `json:"input_path,omitempty"`
		OutputPath  string `json:"output_path,omitempty"`
		KMSProvider string `json:"kms_provider,omitempty"`
		FromKeyID   string `json:"from_kms_key_id,omitempty"`
		ToKeyID     string `json:"to_kms_key_id,omitempty"`
	}
	if err := json.Unmarshal([]byte(C.GoString(cOpts)), &in); err != nil {
		return errReply("invalid JSON: " + err.Error())
	}
	res, err := vaultpack.Rewrap(vaultpack.RewrapOptions{
		InputPath:   in.InputPath,
		OutputPath:  in.OutputPath,
		KMSProvider: in.KMSProvider,
		FromKeyID:   in.FromKeyID,
		ToKeyID:     in.ToKeyID,
	})
	if err != nil {
		return errReply(err.Error())
	}
	return jsonReply(map[string]any{
		"ok":          true,
		"bundle_path": res.OutputPath,
		"manifest":    res.Manifest,
	})
}

func main() {}
