package core

import (
	"net/http"
)

func (t *TLDR) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId, ok := r.Context().Value(requestIdKey).(string)
		if !ok {
			requestId = "NOID"
		}

		t.Logger.Info("Request:", "id", requestId, "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
