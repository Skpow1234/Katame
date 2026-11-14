package crypto

import (
	"testing"
)

// FuzzCombineSharesCorrupted ensures CombineShares / UnmarshalShare reject
// corrupted share data without panicking (M41).
func FuzzCombineSharesCorrupted(f *testing.F) {
	secret := []byte("katame-secret-32-bytes-padding!!")
	shares, err := SplitSecret(secret, 5, 3)
	if err != nil {
		f.Fatal(err)
	}
	wire := MarshalShare(shares[0])
	f.Add(wire)
	f.Add([]byte{})
	f.Add([]byte("KATAME-SHARE-V1\n") )
	f.Add(make([]byte, 64))

	f.Fuzz(func(t *testing.T, data []byte) {
		s, err := UnmarshalShare(data)
		if err != nil {
			return
		}
		// Mix a corrupted share with valid siblings when possible.
		group := []*Share{s}
		for i := 1; i < 3 && i < len(shares); i++ {
			group = append(group, shares[i])
		}
		_, _ = CombineShares(group)
	})
}

// FuzzSplitCombineRoundTrip (M41 companion).
func FuzzSplitCombineRoundTrip(f *testing.F) {
	f.Add([]byte("hello"), byte(5), byte(3))
	f.Add([]byte{1, 2, 3}, byte(3), byte(2))

	f.Fuzz(func(t *testing.T, secret []byte, nByte, kByte byte) {
		if len(secret) == 0 || len(secret) > 1024 {
			return
		}
		n := int(nByte%10) + 2 // 2..11
		k := int(kByte%byte(n-1)) + 2
		if k > n {
			k = n
		}
		shares, err := SplitSecret(secret, n, k)
		if err != nil {
			return
		}
		got, err := CombineShares(shares[:k])
		if err != nil {
			t.Fatalf("combine: %v", err)
		}
		if string(got) != string(secret) {
			t.Fatal("secret mismatch")
		}
	})
}
