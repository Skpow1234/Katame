package serve

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

type authMiddleware struct {
	tokenBytes []byte
	oidc       *oidcValidator
}

func newAuthMiddleware(token string, oidc *oidcValidator) *authMiddleware {
	a := &authMiddleware{oidc: oidc}
	if token != "" {
		a.tokenBytes = []byte(token)
	}
	return a
}

// authenticate returns (principal, ok). Principal is the JWT "sub" when OIDC
// succeeds, "bearer" for a matching static token, or "" when auth is disabled.
func (a *authMiddleware) authenticate(r *http.Request) (string, bool) {
	got := bearerToken(r)
	if got == "" {
		return "", len(a.tokenBytes) == 0 && a.oidc == nil
	}
	if len(a.tokenBytes) > 0 && subtle.ConstantTimeCompare([]byte(got), a.tokenBytes) == 1 {
		return "bearer", true
	}
	if a.oidc != nil {
		sub, err := a.oidc.validate(got)
		if err == nil {
			if sub == "" {
				sub = "oidc"
			}
			return sub, true
		}
	}
	return "", len(a.tokenBytes) == 0 && a.oidc == nil
}

// check is retained for tests and callers that only need a bool.
func (a *authMiddleware) check(r *http.Request) bool {
	_, ok := a.authenticate(r)
	return ok
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if h == "" {
		return ""
	}
	const prefix = "Bearer "
	if len(h) < len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(h[len(prefix):])
}
