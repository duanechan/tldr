package core

import (
	"net/http"
	"time"
)

func (t *TLDR) Health(w http.ResponseWriter, r *http.Request) {
	t.jsonResponse(w, http.StatusOK, HealthResponse{
		Status: "OK",
		Uptime: time.Since(t.startedAt).Round(time.Second).String(),
	})
}
