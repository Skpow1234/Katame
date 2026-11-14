package util

import (
	"regexp"
	"strings"
)

var (
	rePEMBlock = regexp.MustCompile(`(?s)-----BEGIN [^-]*PRIVATE KEY-----.*?-----END [^-]*PRIVATE KEY-----`)
	rePassword = regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[=:]\s*\S+`)
	reLegacy   = regexp.MustCompile(`ed25519-priv:[A-Za-z0-9+/=]+`)
	reBearer   = regexp.MustCompile(`(?i)(bearer\s+)[A-Za-z0-9._\-]+`)
)

const redactPlaceholder = "[REDACTED]"

// RedactSecrets strips private key PEM blocks, password assignments, and
// legacy ed25519-priv material from log/error strings.
func RedactSecrets(s string) string {
	if s == "" {
		return s
	}
	out := rePEMBlock.ReplaceAllString(s, redactPlaceholder)
	out = reLegacy.ReplaceAllString(out, "ed25519-priv:"+redactPlaceholder)
	out = rePassword.ReplaceAllString(out, "${1}="+redactPlaceholder)
	out = reBearer.ReplaceAllString(out, "${1}"+redactPlaceholder)
	// Collapse accidental double equals from password rewrite edge cases.
	out = strings.ReplaceAll(out, "=="+redactPlaceholder, "="+redactPlaceholder)
	return out
}
