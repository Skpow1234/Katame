package serve

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type rateLimiter struct {
	rps   rate.Limit
	burst int
	mu    sync.Mutex
	byKey map[string]*rate.Limiter
}

func newRateLimiter(rps float64, burst int) *rateLimiter {
	if rps <= 0 {
		return nil
	}
	if burst <= 0 {
		burst = int(rps) + 1
	}
	return &rateLimiter{
		rps:   rate.Limit(rps),
		burst: burst,
		byKey: make(map[string]*rate.Limiter),
	}
}

func (rl *rateLimiter) allow(key string) bool {
	if rl == nil {
		return true
	}
	if key == "" {
		key = "unknown"
	}
	rl.mu.Lock()
	lim, ok := rl.byKey[key]
	if !ok {
		lim = rate.NewLimiter(rl.rps, rl.burst)
		rl.byKey[key] = lim
	}
	rl.mu.Unlock()
	return lim.Allow()
}

func rateLimitKey(r *http.Request) string {
	if p := principalFrom(r.Context()); p != "" {
		return "principal:" + p
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "ip:" + r.RemoteAddr
	}
	return "ip:" + ip
}

func (s *Server) rateLimitMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isProbePath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		if s.rateLimit != nil && !s.rateLimit.allow(rateLimitKey(r)) {
			s.metrics.RateLimited.Add(1)
			writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}
		next.ServeHTTP(w, r)
	})
}
