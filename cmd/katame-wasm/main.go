//go:build js && wasm

// Package main builds a WebAssembly module that exposes the Katame SDK
// to JavaScript / TypeScript callers.
//
// Build:
//
//	GOOS=js GOARCH=wasm go build -o katame.wasm ./cmd/katame-wasm
//
// The Go runtime's `wasm_exec.js` is required to instantiate the module.
// After `go.run(instance)`, the module installs a global object
// `globalThis.Katame` with the following methods:
//
//	Katame.version()                         -> string
//	Katame.protect(optsObject)               -> resultObject
//	Katame.decrypt(optsObject)               -> resultObject
//	Katame.inspect(bundleBytes:Uint8Array)   -> resultObject
//	Katame.sign(optsObject)                  -> resultObject (in-memory)
//	Katame.verify(optsObject)                -> resultObject
//
// Every method takes plain JS objects, never file paths — the WASM
// sandbox has no filesystem. Binary blobs are passed as Uint8Array.
package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"syscall/js"

	"github.com/Skpow1234/Katame/pkg/katame"
)

func main() {
	obj := js.Global().Get("Object").New()
	obj.Set("version", js.FuncOf(jsVersion))
	obj.Set("protect", js.FuncOf(jsProtect))
	obj.Set("decrypt", js.FuncOf(jsDecrypt))
	obj.Set("inspect", js.FuncOf(jsInspect))
	obj.Set("sign", js.FuncOf(jsSign))
	obj.Set("verify", js.FuncOf(jsVerify))
	js.Global().Set("Katame", obj)
	select {}
}

func jsError(msg string) any { return map[string]any{"ok": false, "error": msg} }

func bytesFromJS(v js.Value) ([]byte, error) {
	if v.IsUndefined() || v.IsNull() {
		return nil, nil
	}
	if v.Type() == js.TypeString {
		return base64.StdEncoding.DecodeString(v.String())
	}
	if v.InstanceOf(js.Global().Get("Uint8Array")) {
		out := make([]byte, v.Get("length").Int())
		js.CopyBytesToGo(out, v)
		return out, nil
	}
	return nil, errors.New("expected Uint8Array or base64 string")
}

func bytesToJS(b []byte) js.Value {
	u8 := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(u8, b)
	return u8
}

func optString(opts js.Value, key string) string {
	v := opts.Get(key)
	if v.IsUndefined() || v.IsNull() {
		return ""
	}
	return v.String()
}

func optInt(opts js.Value, key string) int {
	v := opts.Get(key)
	if v.IsUndefined() || v.IsNull() {
		return 0
	}
	return v.Int()
}

func manifestJS(m *katame.Manifest) any {
	data, err := katame.MarshalManifest(m)
	if err != nil {
		return nil
	}
	return js.Global().Get("JSON").Call("parse", string(data))
}

func jsVersion(this js.Value, args []js.Value) any { return katame.Version }

func jsProtect(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return jsError("protect: opts object required")
	}
	o := args[0]
	plaintext, err := bytesFromJS(o.Get("plaintext"))
	if err != nil {
		return jsError("protect: " + err.Error())
	}
	if plaintext == nil {
		return jsError("protect: plaintext (Uint8Array) is required")
	}
	key, err := bytesFromJS(o.Get("key"))
	if err != nil {
		return jsError("protect: " + err.Error())
	}
	aad, err := bytesFromJS(o.Get("aad"))
	if err != nil {
		return jsError("protect: " + err.Error())
	}
	opts := katame.ProtectOptions{
		Plaintext:      plaintext,
		Key:            key,
		Password:       optString(o, "password"),
		KDFAlgo:        optString(o, "kdfAlgo"),
		Cipher:         optString(o, "cipher"),
		InputName:      optString(o, "inputName"),
		Compress:       optString(o, "compress"),
		KMSProvider:    optString(o, "kmsProvider"),
		KMSKeyID:       optString(o, "kmsKeyId"),
		HashAlgo:       optString(o, "hashAlgo"),
		AAD:            aad,
		SplitShares:    optInt(o, "splitShares"),
		SplitThreshold: optInt(o, "splitThreshold"),
	}
	if recs := o.Get("recipients"); !recs.IsUndefined() && !recs.IsNull() && recs.Length() > 0 {
		for i := 0; i < recs.Length(); i++ {
			r := recs.Index(i)
			pem := optString(r, "publicKeyPEM")
			path := optString(r, "publicKeyPath")
			rec := katame.Recipient{PublicKeyPath: path}
			if pem != "" {
				rec.PublicKeyPEM = []byte(pem)
			}
			opts.Recipients = append(opts.Recipients, rec)
		}
	}
	var buf bytes.Buffer
	opts.OutputWriter = &buf
	res, err := katame.Protect(opts)
	if err != nil {
		return jsError(err.Error())
	}
	reply := map[string]any{
		"ok":       true,
		"bundle":   bytesToJS(buf.Bytes()),
		"manifest": manifestJS(res.Manifest),
	}
	if res.GeneratedKey != nil {
		reply["generatedKey"] = bytesToJS(res.GeneratedKey)
	}
	if len(res.Shares) > 0 {
		arr := js.Global().Get("Array").New(len(res.Shares))
		for i, sh := range res.Shares {
			item := js.Global().Get("Object").New()
			item.Set("index", sh.Index)
			item.Set("data", bytesToJS(sh.Data))
			arr.SetIndex(i, item)
		}
		reply["shares"] = arr
	}
	return reply
}

func jsDecrypt(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return jsError("decrypt: opts object required")
	}
	o := args[0]
	bundleBytes, err := bytesFromJS(o.Get("bundle"))
	if err != nil || bundleBytes == nil {
		return jsError("decrypt: bundle (Uint8Array) is required")
	}
	key, err := bytesFromJS(o.Get("key"))
	if err != nil {
		return jsError("decrypt: " + err.Error())
	}
	aad, err := bytesFromJS(o.Get("aad"))
	if err != nil {
		return jsError("decrypt: " + err.Error())
	}
	decOpts := katame.DecryptOptions{
		InputBytes:     bundleBytes,
		Key:            key,
		Password:       optString(o, "password"),
		KMSProvider:    optString(o, "kmsProvider"),
		PrivateKeyPath: optString(o, "privateKeyPath"),
		AAD:            aad,
	}
	if pem := optString(o, "privateKeyPEM"); pem != "" {
		decOpts.PrivateKey = []byte(pem)
	}
	res, err := katame.Decrypt(decOpts)
	if err != nil {
		return jsError(err.Error())
	}
	return map[string]any{
		"ok":        true,
		"manifest":  manifestJS(res.Manifest),
		"plaintext": bytesToJS(res.Plaintext),
	}
}

func jsInspect(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return jsError("inspect: bundle (Uint8Array) is required")
	}
	b, err := bytesFromJS(args[0])
	if err != nil || b == nil {
		return jsError("inspect: bundle (Uint8Array) is required")
	}
	m, err := katame.InspectBytes(b)
	if err != nil {
		return jsError(err.Error())
	}
	return map[string]any{"ok": true, "manifest": manifestJS(m)}
}

func jsSign(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return jsError("sign: opts object required")
	}
	o := args[0]
	bundleBytes, err := bytesFromJS(o.Get("bundle"))
	if err != nil || bundleBytes == nil {
		return jsError("sign: bundle (Uint8Array) is required")
	}
	pem := optString(o, "privateKeyPEM")
	if pem == "" {
		return jsError("sign: privateKeyPEM is required")
	}
	signed, res, err := katame.SignBytes(bundleBytes, katame.SignOptions{
		PrivateKey: []byte(pem),
		Algo:       optString(o, "algo"),
	})
	if err != nil {
		return jsError(err.Error())
	}
	return map[string]any{
		"ok":        true,
		"algorithm": res.Algorithm,
		"signedAt":  res.SignedAt,
		"signature": bytesToJS(res.Signature),
		"bundle":    bytesToJS(signed),
	}
}

func jsVerify(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return jsError("verify: opts object required")
	}
	o := args[0]
	bundleBytes, err := bytesFromJS(o.Get("bundle"))
	if err != nil || bundleBytes == nil {
		return jsError("verify: bundle (Uint8Array) is required")
	}
	pem := optString(o, "publicKeyPEM")
	if pem == "" {
		return jsError("verify: publicKeyPEM is required")
	}
	res, err := katame.VerifyBytes(bundleBytes, katame.VerifyOptions{
		PublicKey: []byte(pem),
	})
	if err != nil {
		return jsError(err.Error())
	}
	return map[string]any{
		"ok":        true,
		"valid":     res.Valid,
		"algorithm": res.Algorithm,
		"signedAt":  res.SignedAt,
		"manifest":  manifestJS(res.Manifest),
	}
}
