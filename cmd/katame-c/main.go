//go:build cgo && cshared

// Package main builds a C-shared library (libkatame.so / .dylib / .dll)
// exposing the Katame SDK to native callers (Python, Ruby, Java, etc).
//
// Build:
//
//	go build -buildmode=c-shared -tags cshared -o libkatame.so ./cmd/katame-c
//
// The exported C ABI is intentionally small and JSON-based: every entry
// point takes a single null-terminated JSON string and returns a single
// malloc-allocated JSON string. The caller is responsible for releasing the
// returned string with kt_free.
package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"encoding/json"
	"unsafe"

	_ "github.com/Skpow1234/Katame/internal/kms" // register aws/gcp/azure/mock providers
	"github.com/Skpow1234/Katame/internal/util"
	"github.com/Skpow1234/Katame/pkg/katame"
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

//export kt_version
func kt_version() *C.char {
	return C.CString(katame.Version)
}

//export kt_free
func kt_free(p *C.char) {
	if p != nil {
		C.free(unsafe.Pointer(p))
	}
}

type recipientIn struct {
	PublicKeyPEM  string `json:"public_key_pem,omitempty"`
	PublicKeyPath string `json:"public_key_path,omitempty"`
}

//export kt_protect
func kt_protect(cOpts *C.char) *C.char {
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

	opts := katame.ProtectOptions{
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
		rec := katame.Recipient{PublicKeyPath: r.PublicKeyPath}
		if r.PublicKeyPEM != "" {
			rec.PublicKeyPEM = []byte(r.PublicKeyPEM)
		}
		opts.Recipients = append(opts.Recipients, rec)
	}

	res, err := katame.Protect(opts)
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

//export kt_decrypt
func kt_decrypt(cOpts *C.char) *C.char {
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

	opts := katame.DecryptOptions{
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

	res, err := katame.Decrypt(opts)
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

//export kt_inspect
func kt_inspect(cOpts *C.char) *C.char {
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
			m, err := katame.InspectBytes(b)
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
	m, err := katame.Inspect(path)
	if err != nil {
		return errReply(err.Error())
	}
	return jsonReply(map[string]any{"ok": true, "manifest": m})
}

//export kt_sign
func kt_sign(cOpts *C.char) *C.char {
	var in struct {
		BundlePath     string `json:"bundle_path"`
		PrivateKeyPEM  string `json:"private_key_pem,omitempty"`
		PrivateKeyPath string `json:"private_key_path,omitempty"`
		Algo           string `json:"algo,omitempty"`
	}
	if err := json.Unmarshal([]byte(C.GoString(cOpts)), &in); err != nil {
		return errReply("invalid JSON: " + err.Error())
	}
	opts := katame.SignOptions{
		BundlePath:     in.BundlePath,
		PrivateKeyPath: in.PrivateKeyPath,
		Algo:           in.Algo,
	}
	if in.PrivateKeyPEM != "" {
		opts.PrivateKey = []byte(in.PrivateKeyPEM)
	}
	res, err := katame.SignBundle(opts)
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

//export kt_verify
func kt_verify(cOpts *C.char) *C.char {
	var in struct {
		BundlePath    string `json:"bundle_path"`
		PublicKeyPEM  string `json:"public_key_pem,omitempty"`
		PublicKeyPath string `json:"public_key_path,omitempty"`
	}
	if err := json.Unmarshal([]byte(C.GoString(cOpts)), &in); err != nil {
		return errReply("invalid JSON: " + err.Error())
	}
	opts := katame.VerifyOptions{
		BundlePath:    in.BundlePath,
		PublicKeyPath: in.PublicKeyPath,
	}
	if in.PublicKeyPEM != "" {
		opts.PublicKey = []byte(in.PublicKeyPEM)
	}
	res, err := katame.Verify(opts)
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

//export kt_rewrap
func kt_rewrap(cOpts *C.char) *C.char {
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
	res, err := katame.Rewrap(katame.RewrapOptions{
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
