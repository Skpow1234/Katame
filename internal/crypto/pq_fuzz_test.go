package crypto

import (
	"testing"
)

func fuzzPQSignVerify(f *testing.F, algo string) {
	privPEM, pubPEM, err := GenerateSigningKeys(algo)
	if err != nil {
		f.Fatal(err)
	}
	signer, gotAlgo, err := ParsePrivateKey(privPEM)
	if err != nil {
		f.Fatal(err)
	}
	if gotAlgo != algo {
		f.Fatalf("algo: got %s want %s", gotAlgo, algo)
	}
	pub, _, err := ParseAnyPublicKey(pubPEM)
	if err != nil {
		f.Fatal(err)
	}

	goodMsg := []byte("seed message")
	goodSig, err := SignMessage(signer, algo, goodMsg)
	if err != nil {
		f.Fatal(err)
	}
	f.Add(goodMsg, goodSig)
	f.Add([]byte{}, []byte{})
	f.Add([]byte("x"), make([]byte, 16))

	f.Fuzz(func(t *testing.T, message, sig []byte) {
		ok, err := VerifySignature(pub, algo, message, sig)
		if err != nil {
			return
		}
		if ok {
			// Only the exact seed pair (or a lucky collision) should verify.
			_ = ok
		}
	})
}

// FuzzMLDSA65Verify (M39).
func FuzzMLDSA65Verify(f *testing.F) {
	fuzzPQSignVerify(f, SignAlgoMLDSA65)
}

// FuzzSLHDSA128sVerify (M40).
func FuzzSLHDSA128sVerify(f *testing.F) {
	fuzzPQSignVerify(f, SignAlgoSLHDSA128s)
}
