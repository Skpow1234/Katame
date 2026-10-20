package katame_test

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Skpow1234/Katame/pkg/katame"
)

func TestProtectDecrypt_KeyFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	inPath := filepath.Join(dir, "secret.txt")
	plaintext := []byte("the quick brown fox jumps over the lazy dog")
	if err := os.WriteFile(inPath, plaintext, 0o600); err != nil {
		t.Fatal(err)
	}
	outPath := filepath.Join(dir, "secret.katame")

	res, err := katame.Protect(katame.ProtectOptions{
		InputPath:  inPath,
		OutputPath: outPath,
	})
	if err != nil {
		t.Fatalf("Protect: %v", err)
	}
	if res.GeneratedKey == nil {
		t.Fatal("expected a generated key when no Key/Password supplied")
	}
	if res.Manifest == nil {
		t.Fatal("expected non-nil manifest")
	}
	if res.Manifest.Input.Name != "secret.txt" {
		t.Errorf("input.name = %q, want secret.txt", res.Manifest.Input.Name)
	}
	if res.Manifest.Input.Size != int64(len(plaintext)) {
		t.Errorf("input.size = %d, want %d", res.Manifest.Input.Size, len(plaintext))
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("output bundle not written: %v", err)
	}

	plainOut := filepath.Join(dir, "secret.out")
	dec, err := katame.Decrypt(katame.DecryptOptions{
		InputPath:  outPath,
		OutputPath: plainOut,
		Key:        res.GeneratedKey,
	})
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec.Manifest == nil {
		t.Fatal("expected manifest in DecryptResult")
	}

	got, err := os.ReadFile(plainOut)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("round-trip mismatch: got %q want %q", got, plaintext)
	}
}

func TestProtectDecrypt_InMemory_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	plaintext := []byte("hello, sdk")

	bundlePath := filepath.Join(dir, "mem.katame")
	res, err := katame.Protect(katame.ProtectOptions{
		Plaintext:  plaintext,
		OutputPath: bundlePath,
		InputName:  "mem.txt",
	})
	if err != nil {
		t.Fatalf("Protect: %v", err)
	}
	if res.Manifest.Input.Name != "mem.txt" {
		t.Errorf("input.name = %q, want mem.txt", res.Manifest.Input.Name)
	}

	dec, err := katame.Decrypt(katame.DecryptOptions{
		InputPath: bundlePath,
		Key:       res.GeneratedKey,
	})
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(dec.Plaintext, plaintext) {
		t.Fatalf("plaintext mismatch")
	}
}

func TestProtectDecrypt_Password_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "p.katame")
	plaintext := []byte("super secret")
	password := "correct horse battery staple"

	_, err := katame.Protect(katame.ProtectOptions{
		Plaintext:  plaintext,
		OutputPath: bundlePath,
		Password:   password,
	})
	if err != nil {
		t.Fatalf("Protect: %v", err)
	}

	dec, err := katame.Decrypt(katame.DecryptOptions{
		InputPath: bundlePath,
		Password:  password,
	})
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(dec.Plaintext, plaintext) {
		t.Fatalf("plaintext mismatch")
	}

	if _, err := katame.Decrypt(katame.DecryptOptions{
		InputPath: bundlePath,
		Password:  "wrong password",
	}); err == nil {
		t.Fatal("expected error decrypting with wrong password")
	}
}

func TestDecrypt_WrongKey_Fails(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "wk.katame")
	res, err := katame.Protect(katame.ProtectOptions{
		Plaintext:  []byte("x"),
		OutputPath: bundlePath,
	})
	if err != nil {
		t.Fatal(err)
	}
	wrong := make([]byte, len(res.GeneratedKey))
	for i := range wrong {
		wrong[i] = res.GeneratedKey[i] ^ 0xFF
	}
	if _, err := katame.Decrypt(katame.DecryptOptions{
		InputPath: bundlePath,
		Key:       wrong,
	}); err == nil {
		t.Fatal("expected wrong-key error")
	}
}

func TestProtectDecrypt_InputBytes_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "ib.katame")
	plaintext := []byte("from disk, into memory")

	res, err := katame.Protect(katame.ProtectOptions{
		Plaintext:  plaintext,
		OutputPath: bundlePath,
	})
	if err != nil {
		t.Fatal(err)
	}

	bundleBytes, err := os.ReadFile(bundlePath)
	if err != nil {
		t.Fatal(err)
	}
	dec, err := katame.Decrypt(katame.DecryptOptions{
		InputBytes: bundleBytes,
		Key:        res.GeneratedKey,
	})
	if err != nil {
		t.Fatalf("Decrypt InputBytes: %v", err)
	}
	if !bytes.Equal(dec.Plaintext, plaintext) {
		t.Fatal("plaintext mismatch")
	}
}

func TestInspect(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "i.katame")
	plaintext := []byte("inspectable")
	res, err := katame.Protect(katame.ProtectOptions{
		Plaintext:  plaintext,
		OutputPath: bundlePath,
		Cipher:     "chacha20-poly1305",
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = res

	m, err := katame.Inspect(bundlePath)
	if err != nil {
		t.Fatalf("Inspect: %v", err)
	}
	if m.Encryption.AEAD != "chacha20-poly1305" {
		t.Errorf("AEAD = %q, want chacha20-poly1305", m.Encryption.AEAD)
	}
	if m.Input.Size != int64(len(plaintext)) {
		t.Errorf("size = %d", m.Input.Size)
	}
}

func TestSignVerify_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "sig.katame")
	_, err := katame.Protect(katame.ProtectOptions{
		Plaintext:  []byte("sign me"),
		OutputPath: bundlePath,
	})
	if err != nil {
		t.Fatal(err)
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})
	pubDER, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})

	sres, err := katame.SignBundle(katame.SignOptions{
		BundlePath: bundlePath,
		PrivateKey: privPEM,
	})
	if err != nil {
		t.Fatalf("SignBundle: %v", err)
	}
	if sres.Algorithm != "ed25519" {
		t.Errorf("algorithm = %q, want ed25519", sres.Algorithm)
	}
	if len(sres.Signature) == 0 {
		t.Error("empty signature")
	}

	vres, err := katame.Verify(katame.VerifyOptions{
		BundlePath: bundlePath,
		PublicKey:  pubPEM,
	})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !vres.Valid {
		t.Fatal("signature did not verify")
	}

	// Verify with a different key fails.
	otherPub, _, _ := ed25519.GenerateKey(rand.Reader)
	otherDER, _ := x509.MarshalPKIXPublicKey(otherPub)
	otherPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: otherDER})
	vres2, err := katame.Verify(katame.VerifyOptions{
		BundlePath: bundlePath,
		PublicKey:  otherPEM,
	})
	if err != nil {
		t.Fatalf("Verify with wrong key: %v", err)
	}
	if vres2.Valid {
		t.Fatal("expected signature to fail with wrong public key")
	}
}

func TestProtect_SignInline(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "inline.katame")

	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	privDER, _ := x509.MarshalPKCS8PrivateKey(priv)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})

	res, err := katame.Protect(katame.ProtectOptions{
		Plaintext:  []byte("sign at protect time"),
		OutputPath: bundlePath,
		Sign:       &katame.SignParams{PrivateKey: privPEM},
	})
	if err != nil {
		t.Fatalf("Protect with Sign: %v", err)
	}
	if res.SignatureAlgo != "ed25519" {
		t.Errorf("signature algo = %q", res.SignatureAlgo)
	}
	if len(res.Signature) == 0 {
		t.Error("expected signature in result")
	}
}

func TestProtect_ErrorPaths(t *testing.T) {
	tests := []struct {
		name string
		opts katame.ProtectOptions
		msg  string
	}{
		{
			name: "no input",
			opts: katame.ProtectOptions{OutputPath: "/tmp/x.katame"},
			msg:  "one of",
		},
		{
			name: "no output",
			opts: katame.ProtectOptions{Plaintext: []byte("x")},
			msg:  "OutputPath or OutputWriter",
		},
		{
			name: "key+password",
			opts: katame.ProtectOptions{
				Plaintext:  []byte("x"),
				OutputPath: filepath.Join(t.TempDir(), "x.katame"),
				Key:        make([]byte, 32),
				Password:   "p",
			},
			msg: "mutually exclusive",
		},
		{
			name: "wrong key length",
			opts: katame.ProtectOptions{
				Plaintext:  []byte("x"),
				OutputPath: filepath.Join(t.TempDir(), "x.katame"),
				Key:        []byte("short"),
			},
			msg: "Key must be 32 bytes",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := katame.Protect(tc.opts)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tc.msg) {
				t.Errorf("error = %v; want substring %q", err, tc.msg)
			}
		})
	}
}

func TestDecrypt_ErrorPaths(t *testing.T) {
	_, err := katame.Decrypt(katame.DecryptOptions{InputPath: "no.katame", Key: make([]byte, 32)})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "k.katame")
	res, err := katame.Protect(katame.ProtectOptions{
		Plaintext:  []byte("x"),
		OutputPath: bundlePath,
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = res
	_, err = katame.Decrypt(katame.DecryptOptions{InputPath: bundlePath})
	if err == nil || !strings.Contains(err.Error(), "Key") {
		t.Fatalf("expected no-key error, got %v", err)
	}
}

func TestVersion(t *testing.T) {
	if katame.Version == "" {
		t.Fatal("Version must not be empty")
	}
	parts := strings.Split(katame.Version, ".")
	if len(parts) != 3 {
		t.Fatalf("Version must be semver, got %q", katame.Version)
	}
}

func TestProtectDecrypt_Compression_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "gzip.katame")
	plaintext := bytes.Repeat([]byte("compressible payload "), 100)

	res, err := katame.Protect(katame.ProtectOptions{
		Plaintext:  plaintext,
		OutputPath: bundlePath,
		Compress:   "gzip",
	})
	if err != nil {
		t.Fatalf("Protect: %v", err)
	}

	dec, err := katame.Decrypt(katame.DecryptOptions{
		InputPath: bundlePath,
		Key:       res.GeneratedKey,
	})
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(dec.Plaintext, plaintext) {
		t.Fatal("plaintext mismatch after compression round-trip")
	}
}

func TestProtectDecrypt_KMS_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "kms.katame")
	plaintext := []byte("kms-wrapped dek")

	_, err := katame.Protect(katame.ProtectOptions{
		Plaintext:   plaintext,
		OutputPath:  bundlePath,
		KMSProvider: "mock",
		KMSKeyID:    "mock-key-id",
	})
	if err != nil {
		t.Fatalf("Protect: %v", err)
	}

	dec, err := katame.Decrypt(katame.DecryptOptions{
		InputPath:   bundlePath,
		KMSProvider: "mock",
	})
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(dec.Plaintext, plaintext) {
		t.Fatal("plaintext mismatch after KMS round-trip")
	}
}

func TestProtectDecrypt_SplitShares_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "split.katame")
	plaintext := []byte("shamir split key material")

	res, err := katame.Protect(katame.ProtectOptions{
		Plaintext:      plaintext,
		OutputPath:     bundlePath,
		SplitShares:    5,
		SplitThreshold: 3,
	})
	if err != nil {
		t.Fatalf("Protect: %v", err)
	}
	if len(res.Shares) != 5 {
		t.Fatalf("expected 5 shares, got %d", len(res.Shares))
	}

	rawShares := [][]byte{res.Shares[0].Data, res.Shares[2].Data, res.Shares[4].Data}
	key, err := katame.CombineShares(rawShares)
	if err != nil {
		t.Fatalf("CombineShares: %v", err)
	}

	dec, err := katame.Decrypt(katame.DecryptOptions{
		InputPath: bundlePath,
		Key:       key,
	})
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(dec.Plaintext, plaintext) {
		t.Fatal("plaintext mismatch after split-share round-trip")
	}
}

func TestRewrap_KMS(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "rewrap.katame")
	outPath := filepath.Join(dir, "rewrap-out.katame")

	_, err := katame.Protect(katame.ProtectOptions{
		Plaintext:   []byte("rewrap me"),
		OutputPath:  bundlePath,
		KMSProvider: "mock",
		KMSKeyID:    "mock-key-id",
	})
	if err != nil {
		t.Fatalf("Protect: %v", err)
	}

	_, err = katame.Rewrap(katame.RewrapOptions{
		InputPath:   bundlePath,
		OutputPath:  outPath,
		KMSProvider: "mock",
		FromKeyID:   "mock-key-id",
		ToKeyID:     "mock-key-id",
	})
	if err != nil {
		t.Fatalf("Rewrap: %v", err)
	}

	dec, err := katame.Decrypt(katame.DecryptOptions{
		InputPath:   outPath,
		KMSProvider: "mock",
	})
	if err != nil {
		t.Fatalf("Decrypt after rewrap: %v", err)
	}
	if string(dec.Plaintext) != "rewrap me" {
		t.Fatalf("plaintext = %q", dec.Plaintext)
	}
}