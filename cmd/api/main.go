package main

import (
	"log"
	"net/http"

	"github.com/duanechan/tldr/internal/core"
	_ "modernc.org/sqlite"
)

func main() {
	app, err := core.New()
	if err != nil {
		log.Fatalf("Failed to initialize app: %s", err.Error())
	}
	defer app.CloseDB()

	app.Use(app.PanicRecoveryMiddleware, core.CorsMiddleware, app.LogMiddleware)

	api := http.NewServeMux()
	api.HandleFunc("POST /v1/auth/register", app.Register)
	api.HandleFunc("POST /v1/auth/login", app.Login)
	api.Handle("POST /v1/summarize/file", app.AuthMiddleware(http.HandlerFunc(app.SummarizeFile)))
	api.Handle("POST /v1/summarize/text", app.AuthMiddleware(http.HandlerFunc(app.SummarizeText)))

	app.Handle("/", http.FileServer(http.Dir("web/dist")))
	app.Handle("/api/", http.StripPrefix("/api", api))
	app.HandleFunc("GET /health", app.Health)

	app.Logger.Info("Server started:", "port", app.Config.Port, "environment", app.Config.Environment)
	if err := http.ListenAndServe(":"+app.Config.Port, app.Handler); err != nil {
		app.Logger.Error("Error occurred:", "error", err.Error())
	}
}
