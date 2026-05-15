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

	app.Use(
		app.PanicRecoveryMiddleware,
		core.CorsMiddleware,
		app.LogMiddleware,
		app.RequestIDMiddleware,
	)

	api := http.NewServeMux()
	api.Handle("POST /v1/auth/register", app.RateLimiterMiddleware(http.HandlerFunc(app.Register)))
	api.Handle("POST /v1/auth/login", app.RateLimiterMiddleware(http.HandlerFunc(app.Login)))
	api.Handle("POST /v1/auth/refresh", http.HandlerFunc(app.Refresh))
	api.Handle("POST /v1/auth/logout", http.HandlerFunc(app.Logout))

	api.Handle("POST /v1/summarize/file", app.AuthMiddleware(http.HandlerFunc(app.SummarizeFile)))
	api.Handle("POST /v1/summarize/text", app.AuthMiddleware(http.HandlerFunc(app.SummarizeText)))

	api.Handle("GET /v1/tldrs", app.AuthMiddleware(http.HandlerFunc(app.UserGetTLDRs)))
	api.Handle("GET /v1/tldrs/{id}", app.AuthMiddleware(http.HandlerFunc(app.UserGetTLDR)))
	api.Handle("PATCH /v1/tldrs/{id}", app.AuthMiddleware(http.HandlerFunc(app.UserUpdateTLDR)))
	api.Handle("DELETE /v1/tldrs/{id}", app.AuthMiddleware(http.HandlerFunc(app.UserDeleteTLDR)))

	api.Handle("GET /v1/me", app.AuthMiddleware(http.HandlerFunc(app.UserGetMe)))

	admin := http.NewServeMux()
	admin.Handle("GET /tldrs", http.HandlerFunc(app.AdminGetTLDRs))
	admin.Handle("GET /tldrs/{id}", http.HandlerFunc(app.AdminGetTLDR))
	admin.Handle("PATCH /tldrs/{id}", http.HandlerFunc(app.AdminUpdateTLDR))
	admin.Handle("DELETE /tldrs/{id}", http.HandlerFunc(app.AdminDeleteTLDR))

	admin.Handle("GET /users", http.HandlerFunc(app.AdminGetUsers))
	admin.Handle("GET /users/{id}", http.HandlerFunc(app.AdminGetUser))
	admin.Handle("PATCH /users/{id}", http.HandlerFunc(app.AdminUpdateUsername))
	admin.Handle("PATCH /users/{id}/reset-password", http.HandlerFunc(app.AdminUpdatePassword))
	admin.Handle("DELETE /users/{id}", http.HandlerFunc(app.AdminDeleteUser))

	app.Handle("/", http.FileServer(http.Dir("web/dist")))
	app.Handle("/admin/", http.StripPrefix("/admin", app.AuthMiddleware(app.AdminMiddleware(admin))))
	app.Handle("/api/", http.StripPrefix("/api", api))
	app.Handle("PATCH /change-username", app.AuthMiddleware(http.HandlerFunc(app.UserUpdateUsername)))
	app.Handle("PATCH /reset-password", app.AuthMiddleware(http.HandlerFunc(app.UserUpdatePassword)))
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
