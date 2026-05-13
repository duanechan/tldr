package core

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (t *TLDR) RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.Must(uuid.NewRandom())
		ctx := context.WithValue(r.Context(), requestIdKey, requestId.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
