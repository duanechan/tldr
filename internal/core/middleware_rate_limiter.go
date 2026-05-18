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

func (a *App) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			a.errorResponse(
				w,
				r.Context(),
				http.StatusBadRequest,
				"Invalid IP address",
			)
			return
		}

		a.mu.Lock()
		limiter, exists := a.clients[host]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(rateLimitInterval), burstLimit)
			a.clients[host] = limiter
		}
		a.mu.Unlock()

		if !limiter.Allow() {
			a.errorResponse(
				w,
				r.Context(),
				http.StatusTooManyRequests,
				"Rate limit exceeded, try again later",
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
