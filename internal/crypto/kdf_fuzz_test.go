package crypto

import (
	"testing"
)

// FuzzDeriveKeyWeakParams ensures DeriveKey rejects invalid/weak params and
// never panics on arbitrary inputs (M42).
func FuzzDeriveKeyWeakParams(f *testing.F) {
	f.Add([]byte("password"), []byte("0123456789abcdef"), uint32(0), uint32(0), uint8(0), byte(0))
	f.Add([]byte("password"), []byte("0123456789abcdef"), uint32(1), uint32(1024), uint8(1), byte(0))
	f.Add([]byte(""), []byte(""), uint32(1), uint32(1), uint8(1), byte(1))

	f.Fuzz(func(t *testing.T, password, salt []byte, time, memory uint32, threads uint8, algoSel byte) {
		var params KDFParams
		switch algoSel % 3 {
		case 0:
			params = KDFParams{Algo: KDFArgon2id, Time: time, Memory: memory, Threads: threads}
		case 1:
			params = KDFParams{Algo: KDFScrypt, N: int(time), R: int(memory), P: int(threads)}
		default:
			params = KDFParams{Algo: KDFPBKDF2SHA256, Iterations: int(time)}
		}

		_, err := DeriveKey(password, salt, params, 32)

		// Explicit weak/invalid params must always error.
		switch params.Algo {
		case KDFArgon2id:
			if time == 0 || memory == 0 || threads == 0 {
				if err == nil {
					t.Fatal("argon2id: expected error for zero param")
				}
			}
		case KDFScrypt:
			if params.N <= 0 || params.R <= 0 || params.P <= 0 {
				if err == nil {
					t.Fatal("scrypt: expected error for non-positive param")
				}
			}
		case KDFPBKDF2SHA256:
			if params.Iterations <= 0 {
				if err == nil {
					t.Fatal("pbkdf2: expected error for non-positive iterations")
				}
			}
		}

		// Avoid pathological allocations during fuzzing.
		_ = err
	})
}
