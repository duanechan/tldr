package main

import (
	"net/http"

	"github.com/duanechan/tldr/internal/core"
)

func main() {
	app, err := core.New()
	if err != nil {
		app.Logger.Error("Failed to initialize app:", "error", err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir("web/dist")))
	mux.HandleFunc("GET /health", app.Health)
	mux.HandleFunc("POST /api/v1/summarize/document", app.SummarizeDocument)

	app.Logger.Info("Server started:", "port", app.Config.Port, "environment", app.Config.Environment)
	if err := http.ListenAndServe(":"+app.Config.Port, app.LogMiddleware(mux)); err != nil {
		app.Logger.Error("Error occured:", "error", err.Error())
	}
}
