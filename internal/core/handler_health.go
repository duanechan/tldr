package core

import (
	"net/http"
	"time"
)

type services struct {
	Database string `json:"database"`
	Model    string `json:"model"`
}

type healthResponse struct {
	Status   string   `json:"status"`
	Uptime   string   `json:"uptime"`
	Services services `json:"services"`
}

func (a *App) Health(w http.ResponseWriter, r *http.Request) {
	dbConnected := "UP"
	if err := a.db.Ping(); err != nil {
		dbConnected = "DOWN"
	}

	modelConnected := "UP"
	if _, err := a.Client.Models.List(r.Context(), nil); err != nil {
		modelConnected = "DOWN"
	}

	a.jsonResponse(w, http.StatusOK, healthResponse{
		Status: "OK",
		Uptime: time.Since(a.startedAt).Round(time.Second).String(),
		Services: services{
			Database: dbConnected,
			Model:    modelConnected,
		},
	})
}
