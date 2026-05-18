package core

import (
	"net/http"
	"time"
)

func (t *TLDR) Health(w http.ResponseWriter, r *http.Request) {
	dbConnected := "UP"
	if err := t.db.Ping(); err != nil {
		dbConnected = "DOWN"
	}

	modelConnected := "UP"
	if _, err := t.Client.Models.List(r.Context(), nil); err != nil {
		modelConnected = "DOWN"
	}

	t.jsonResponse(w, http.StatusOK, HealthResponse{
		Status: "OK",
		Uptime: time.Since(t.startedAt).Round(time.Second).String(),
		Services: Services{
			Database: dbConnected,
			Model:    modelConnected,
		},
	})
}
