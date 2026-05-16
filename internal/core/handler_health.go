package core

import (
	"net/http"
	"time"

	"github.com/duanechan/tldr/internal/types"
)

func (t *TLDR) Health(w http.ResponseWriter, r *http.Request) {
	t.jsonResponse(w, http.StatusOK, types.HealthResponse{
		Status: "OK",
		Uptime: time.Since(t.startedAt).Round(time.Second).String(),
	})
}
