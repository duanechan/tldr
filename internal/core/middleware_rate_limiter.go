package core

import (
	"net"
	"net/http"

	"github.com/duanechan/tldr/internal/types"
	"golang.org/x/time/rate"
)

const ()

func (t *TLDR) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			t.errorResponse(w, r.Context(), http.StatusBadRequest, "Invalid IP address")
			return
		}

		t.mu.Lock()
		limiter, exists := t.clients[host]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(types.RateLimitInterval), types.BurstLimit)
			t.clients[host] = limiter
		}
		t.mu.Unlock()

		if !limiter.Allow() {
			t.errorResponse(w, r.Context(), http.StatusTooManyRequests, "Rate limit exceeded, try again later")
			return
		}

		next.ServeHTTP(w, r)
	})
}
