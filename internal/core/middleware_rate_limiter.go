package core

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const (
	rateLimitInterval = 1 * time.Minute
	burstLimit        = 5
)

func (t *TLDR) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.mu.Lock()
		defer t.mu.Unlock()
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			t.errorResponse(w, r.Context(), http.StatusBadRequest, "Invalid IP address")
			return
		}

		limiter, exists := t.clients[host]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(rateLimitInterval), burstLimit)
			t.clients[host] = limiter
		}

		if !limiter.Allow() {
			t.errorResponse(w, r.Context(), http.StatusTooManyRequests, "Rate limit exceeded, try again later")
			return
		}

		next.ServeHTTP(w, r)
	})
}
