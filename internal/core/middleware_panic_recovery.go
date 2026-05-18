package core

import (
	"net/http"
	"runtime/debug"
)

func (a *App) PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				a.Logger.Error(
					"Server panicked:",
					"stack",
					string(debug.Stack()),
					"value",
					v,
				)
				a.errorResponse(
					w,
					r.Context(),
					http.StatusInternalServerError,
					"Something went wrong",
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
