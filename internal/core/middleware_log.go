package core

import (
	"net/http"
)

func (t *TLDR) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId, _ := r.Context().Value(requestIdKey).(string)

		t.Logger.Info("Request:", "id", requestId, "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
