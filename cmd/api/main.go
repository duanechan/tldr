package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	api.HandleFunc("POST /v1/auth/refresh", app.Refresh)

	api.Handle("POST /v1/summarize/file", app.AuthMiddleware(http.HandlerFunc(app.SummarizeFile)))
	api.Handle("POST /v1/summarize/text", app.AuthMiddleware(http.HandlerFunc(app.SummarizeText)))

	api.Handle("GET /v1/tldrs", app.AuthMiddleware(http.HandlerFunc(app.GetTLDRs)))
	api.Handle("GET /v1/tldrs/{id}", app.AuthMiddleware(http.HandlerFunc(app.GetTLDR)))
	api.Handle("PUT /v1/tldrs/{id}", app.AuthMiddleware(http.HandlerFunc(app.UpdateTLDR)))
	api.Handle("DELETE /v1/tldrs/{id}", app.AuthMiddleware(http.HandlerFunc(app.DeleteTLDR)))

	app.Handle("/", http.FileServer(http.Dir("web/dist")))
	app.Handle("/api/", http.StripPrefix("/api", api))
	app.HandleFunc("GET /health", app.Health)

	server := &http.Server{
		Addr:    ":" + app.Config.Port,
		Handler: app.Handler,
	}

	app.Logger.Info("Server started:", "port", app.Config.Port, "environment", app.Config.Environment)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Error("Error occurred:", "error", err.Error())
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	app.Logger.Info("Server shutting down...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		app.Logger.Error("Error occurred shutting down:", "error", err.Error())
	}
}
