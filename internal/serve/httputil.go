package serve

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type ctxKey int

const (
	ctxKeyRequestID ctxKey = iota
	ctxKeyPrincipal
)

// requestIDFrom returns the request ID attached by requestIDMW.
func requestIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeyRequestID).(string)
	return v
}

// principalFrom returns the authenticated principal (JWT sub or "bearer").
func principalFrom(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeyPrincipal).(string)
	return v
}

func (s *Server) requestIDMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			id = newRequestID()
		}
		w.Header().Set("X-Request-Id", id)
		ctx := context.WithValue(r.Context(), ctxKeyRequestID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func newRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "req-unknown"
	}
	return hex.EncodeToString(b[:])
}

func errorCodeForStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "policy_denied"
	case http.StatusMethodNotAllowed:
		return "method_not_allowed"
	case http.StatusTooManyRequests:
		return "rate_limited"
	case http.StatusRequestEntityTooLarge:
		return "payload_too_large"
	case http.StatusServiceUnavailable:
		return "unavailable"
	case http.StatusInternalServerError:
		return "internal_error"
	default:
		if status >= 400 && status < 500 {
			return "client_error"
		}
		return "internal_error"
	}
}

func isProbePath(path string) bool {
	return path == "/healthz" || path == "/readyz" || path == "/metrics"
}
