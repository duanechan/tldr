package core

import (
	"net/http"
	"time"
)

func (a *App) Health(w http.ResponseWriter, r *http.Request) {
	dbConnected := "UP"
	if err := a.db.Ping(); err != nil {
		dbConnected = "DOWN"
	}

	modelConnected := "UP"
	if _, err := a.Client.Models.List(r.Context(), nil); err != nil {
		modelConnected = "DOWN"
	}

	a.jsonResponse(w, http.StatusOK, HealthResponse{
		Status: "OK",
		Uptime: time.Since(a.startedAt).Round(time.Second).String(),
		Services: Services{
			Database: dbConnected,
			Model:    modelConnected,
		},
	})
}
