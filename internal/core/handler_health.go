package core

import (
	"net/http"
	"time"
)

type HealthResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

func (t *TLDR) Health(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, HealthResponse{
		Status: "OK",
		Uptime: time.Since(t.startedAt).Round(time.Second).String(),
	})
}
