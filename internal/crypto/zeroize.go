package crypto

// Zeroize overwrites b with zeros. Use after AEAD operations when a
// temporary DEK or key material is no longer needed.
func Zeroize(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
