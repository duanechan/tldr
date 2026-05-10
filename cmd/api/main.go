package main

import (
	"log"
	"net/http"

	"github.com/duanechan/tldr/internal/core"
)

func main() {
	app, err := core.New()
	if err != nil {
		log.Fatalf("Failed to initialize app: %s", err.Error())
	}

	api := http.NewServeMux()
	api.HandleFunc("POST /v1/summarize/document", app.SummarizeDocument)
	api.HandleFunc("POST /v1/summarize/text", app.SummarizeText)

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir("web/dist")))
	mux.Handle("GET /api/", http.StripPrefix("/api", api))
	mux.HandleFunc("GET /health", app.Health)

	app.Logger.Info("Server started:", "port", app.Config.Port, "environment", app.Config.Environment)
	if err := http.ListenAndServe(":"+app.Config.Port, app.LogMiddleware(mux)); err != nil {
		app.Logger.Error("Error occurred:", "error", err.Error())
	}
}
