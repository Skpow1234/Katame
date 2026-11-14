package crypto

import (
	"bytes"
	"crypto/rand"
	"testing"
)

// TestStreamLastChunkTruncationAllCiphers verifies truncation is detected for
// every AEAD (M44).
func TestStreamLastChunkTruncationAllCiphers(t *testing.T) {
	chunkSize := 64
	for _, cipherName := range SupportedCiphers {
		t.Run(cipherName, func(t *testing.T) {
			info, err := GetCipherInfo(cipherName)
			if err != nil {
				t.Fatal(err)
			}
			plaintext := make([]byte, chunkSize*3)
			rand.Read(plaintext)
			key, _ := GenerateKey(AES256KeySize)

			var cipherBuf bytes.Buffer
			result, err := EncryptStream(bytes.NewReader(plaintext), &cipherBuf, key, nil, chunkSize, cipherName)
			if err != nil {
				t.Fatalf("encrypt: %v", err)
			}

			full := cipherBuf.Bytes()
			// Drop the final encrypted chunk (plaintext chunk + tag).
			cut := len(full) - (chunkSize + info.TagSize)
			if cut <= 0 {
				t.Fatal("ciphertext too short to truncate")
			}
			truncated := full[:cut]

			var out bytes.Buffer
			err = DecryptStream(bytes.NewReader(truncated), &out, key, result.BaseNonce, nil, chunkSize, cipherName)
			if err == nil {
				t.Fatal("expected truncation error")
			}
		})
	}
}

// TestStreamReorderDuplicateChunks verifies reordered or duplicated chunks
// fail authentication (M45).
func TestStreamReorderDuplicateChunks(t *testing.T) {
	chunkSize := 32
	for _, cipherName := range SupportedCiphers {
		t.Run(cipherName, func(t *testing.T) {
			info, err := GetCipherInfo(cipherName)
			if err != nil {
				t.Fatal(err)
			}
			encChunk := chunkSize + info.TagSize
			plaintext := make([]byte, chunkSize*3)
			rand.Read(plaintext)
			key, _ := GenerateKey(AES256KeySize)

			var cipherBuf bytes.Buffer
			result, err := EncryptStream(bytes.NewReader(plaintext), &cipherBuf, key, nil, chunkSize, cipherName)
			if err != nil {
				t.Fatalf("encrypt: %v", err)
			}
			full := cipherBuf.Bytes()
			if len(full) < encChunk*3 {
				t.Fatalf("unexpected ciphertext length %d", len(full))
			}

			chunks := [][]byte{
				append([]byte(nil), full[0:encChunk]...),
				append([]byte(nil), full[encChunk:2*encChunk]...),
				append([]byte(nil), full[2*encChunk:]...),
			}

			t.Run("swap-first-two", func(t *testing.T) {
				reordered := append(append([]byte(nil), chunks[1]...), chunks[0]...)
				reordered = append(reordered, chunks[2]...)
				var out bytes.Buffer
				if err := DecryptStream(bytes.NewReader(reordered), &out, key, result.BaseNonce, nil, chunkSize, cipherName); err == nil {
					t.Fatal("expected error for reordered chunks")
				}
			})

			t.Run("duplicate-first", func(t *testing.T) {
				dup := append(append([]byte(nil), chunks[0]...), chunks[0]...)
				dup = append(dup, chunks[2]...)
				var out bytes.Buffer
				if err := DecryptStream(bytes.NewReader(dup), &out, key, result.BaseNonce, nil, chunkSize, cipherName); err == nil {
					t.Fatal("expected error for duplicated chunk")
				}
			})
		})
	}
}

// TestStreamChunkSizeHardCap (M43).
func TestStreamChunkSizeHardCap(t *testing.T) {
	key, _ := GenerateKey(AES256KeySize)
	tooBig := MaxChunkSize + 1
	_, err := EncryptStream(bytes.NewReader([]byte("x")), &bytes.Buffer{}, key, nil, tooBig, CipherAES256GCM)
	if err == nil {
		t.Fatal("expected error for oversized chunk")
	}
	err = DecryptStream(bytes.NewReader([]byte("x")), &bytes.Buffer{}, key, make([]byte, 12), nil, tooBig, CipherAES256GCM)
	if err == nil {
		t.Fatal("expected error for oversized decrypt chunk")
	}
}

// TestEncryptStreamParallelRace hammers parallel encrypt/decrypt under the race
// detector (M46). Run with: go test -race -run TestEncryptStreamParallelRace
func TestEncryptStreamParallelRace(t *testing.T) {
	key, _ := GenerateKey(AES256KeySize)
	pt := make([]byte, 8*DefaultChunkSize)
	rand.Read(pt)

	const workers = 8
	var ct bytes.Buffer
	res, err := EncryptStreamParallel(bytes.NewReader(pt), &ct, key, []byte("aad"), DefaultChunkSize, CipherAES256GCM, workers)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	var out bytes.Buffer
	if err := DecryptStreamParallel(bytes.NewReader(ct.Bytes()), &out, key, res.BaseNonce, []byte("aad"), DefaultChunkSize, CipherAES256GCM, workers); err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if !bytes.Equal(out.Bytes(), pt) {
		t.Fatal("plaintext mismatch")
	}
}
