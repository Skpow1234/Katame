package keysource

// PlatformSchemeStatus reports whether a reserved OS/HSM scheme is
// implemented in this binary build.
type PlatformSchemeStatus struct {
	Scheme      string `json:"scheme"`
	Implemented bool   `json:"implemented"`
	Details     string `json:"details,omitempty"`
}

// PlatformSchemes returns the availability of hardware/OS key backends.
func PlatformSchemes() []PlatformSchemeStatus {
	return []PlatformSchemeStatus{
		{Scheme: "dpapi", Implemented: dpapiAvailable, Details: dpapiDetails},
		{Scheme: "keychain", Implemented: keychainAvailable, Details: keychainDetails},
		{Scheme: "pkcs11", Implemented: pkcs11Available, Details: pkcs11Details},
		{Scheme: "piv", Implemented: pivAvailable, Details: pivDetails},
	}
}
