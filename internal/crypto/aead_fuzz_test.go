package crypto

import (
	"testing"
)

// FuzzDecryptAESGCM ensures DecryptAESGCM handles arbitrary ciphertext, nonce,
// tag, and AAD without panicking.
func FuzzDecryptAESGCM(f *testing.F) {
	key := make([]byte, AES256KeySize)
	for i := range key {
		key[i] = byte(i)
	}
	res, _ := EncryptAESGCM([]byte("hello world"), key, nil)
	if res != nil {
		f.Add(res.Ciphertext, key[:], res.Nonce, res.Tag, []byte{})
	}
	f.Add(make([]byte, 32), make([]byte, 32), make([]byte, 12), make([]byte, 16), []byte{})

	f.Fuzz(func(t *testing.T, ciphertext, key, nonce, tag, aad []byte) {
		_, _ = DecryptAESGCM(ciphertext, key, nonce, tag, aad)
	})
}

// FuzzEncryptDecryptRoundTrip checks that encrypt→decrypt always round-trips
// when given valid parameters.
func FuzzEncryptDecryptRoundTrip(f *testing.F) {
	f.Add([]byte("hello"), []byte("extra-aad"))
	f.Add([]byte{}, []byte{})
	f.Add(make([]byte, 1024), []byte("large"))

	f.Fuzz(func(t *testing.T, plaintext, aad []byte) {
		key, err := GenerateKey(AES256KeySize)
		if err != nil {
			t.Fatal(err)
		}

		res, err := EncryptAESGCM(plaintext, key, aad)
		if err != nil {
			t.Fatalf("encrypt: %v", err)
		}

		got, err := DecryptAESGCM(res.Ciphertext, key, res.Nonce, res.Tag, aad)
		if err != nil {
			t.Fatalf("decrypt: %v", err)
		}

		if len(plaintext) == 0 && len(got) == 0 {
			return
		}
		if string(got) != string(plaintext) {
			t.Errorf("round-trip mismatch: got %q, want %q", got, plaintext)
		}
	})
}

func fuzzDecryptAEAD(f *testing.F, cipherName string) {
	key := make([]byte, AES256KeySize)
	for i := range key {
		key[i] = byte(i + 1)
	}
	info, err := GetCipherInfo(cipherName)
	if err != nil {
		f.Fatal(err)
	}
	res, err := EncryptAEAD(cipherName, []byte("seed"), key, []byte("aad"))
	if err == nil && res != nil {
		f.Add(res.Ciphertext, append([]byte(nil), key...), res.Nonce, res.Tag, []byte("aad"))
	}
	f.Add(make([]byte, 16), make([]byte, 32), make([]byte, info.NonceSize), make([]byte, info.TagSize), []byte{})

	f.Fuzz(func(t *testing.T, ciphertext, key, nonce, tag, aad []byte) {
		_, _ = DecryptAEAD(cipherName, ciphertext, key, nonce, tag, aad)
	})
}

func fuzzEncryptDecryptAEADRoundTrip(f *testing.F, cipherName string) {
	f.Add([]byte("hello"), []byte("aad"))
	f.Add([]byte{}, []byte{})
	f.Add(make([]byte, 512), []byte("x"))

	f.Fuzz(func(t *testing.T, plaintext, aad []byte) {
		key, err := GenerateKey(AES256KeySize)
		if err != nil {
			t.Fatal(err)
		}
		res, err := EncryptAEAD(cipherName, plaintext, key, aad)
		if err != nil {
			t.Fatalf("encrypt: %v", err)
		}
		got, err := DecryptAEAD(cipherName, res.Ciphertext, key, res.Nonce, res.Tag, aad)
		if err != nil {
			t.Fatalf("decrypt: %v", err)
		}
		if string(got) != string(plaintext) {
			t.Errorf("round-trip mismatch")
		}
	})
}

// FuzzDecryptChaCha20Poly1305 ensures DecryptAEAD never panics (M31).
func FuzzDecryptChaCha20Poly1305(f *testing.F) {
	fuzzDecryptAEAD(f, CipherChaCha20Poly1305)
}

// FuzzEncryptDecryptChaCha20RoundTrip (M31).
func FuzzEncryptDecryptChaCha20RoundTrip(f *testing.F) {
	fuzzEncryptDecryptAEADRoundTrip(f, CipherChaCha20Poly1305)
}

// FuzzDecryptXChaCha20Poly1305 ensures DecryptAEAD never panics (M32).
func FuzzDecryptXChaCha20Poly1305(f *testing.F) {
	fuzzDecryptAEAD(f, CipherXChaCha20Poly1305)
}

// FuzzEncryptDecryptXChaCha20RoundTrip (M32).
func FuzzEncryptDecryptXChaCha20RoundTrip(f *testing.F) {
	fuzzEncryptDecryptAEADRoundTrip(f, CipherXChaCha20Poly1305)
}

// FuzzDecryptAES256GCMSIV ensures DecryptAEAD never panics (M33).
func FuzzDecryptAES256GCMSIV(f *testing.F) {
	fuzzDecryptAEAD(f, CipherAES256GCMSIV)
}

// FuzzEncryptDecryptAES256GCMSIVRoundTrip (M33).
func FuzzEncryptDecryptAES256GCMSIVRoundTrip(f *testing.F) {
	fuzzEncryptDecryptAEADRoundTrip(f, CipherAES256GCMSIV)
}
