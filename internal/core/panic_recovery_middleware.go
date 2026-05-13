package core

import (
	"net/http"
	"runtime/debug"
)

func (t *TLDR) PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				t.Logger.Error("Server panicked:", "stack", string(debug.Stack()), "value", v)
				t.errorResponse(w, http.StatusInternalServerError, "Something went wrong")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
