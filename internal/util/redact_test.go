package util

import (
	"strings"
	"testing"
)

func TestRedactSecrets(t *testing.T) {
	pem := "before -----BEGIN PRIVATE KEY-----\nABC\n-----END PRIVATE KEY----- after"
	got := RedactSecrets(pem)
	if strings.Contains(got, "BEGIN PRIVATE KEY") || !strings.Contains(got, "[REDACTED]") {
		t.Fatalf("pem not redacted: %q", got)
	}
	pw := "password=supersecret"
	got = RedactSecrets(pw)
	if strings.Contains(got, "supersecret") {
		t.Fatalf("password not redacted: %q", got)
	}
	legacy := "ed25519-priv:YWJjZGVm"
	got = RedactSecrets(legacy)
	if strings.Contains(got, "YWJjZGVm") {
		t.Fatalf("legacy key not redacted: %q", got)
	}
}
