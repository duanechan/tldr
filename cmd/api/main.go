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
	dep, err := core.Bootstrap()
	if err != nil {
		log.Fatalf("Failed to build dependencies app: %s", err.Error())
	}

	app := core.New(dep)
	defer app.CloseDB()

	app.Use(
		app.PanicRecoveryMiddleware,
		core.CorsMiddleware,
		app.LogMiddleware,
		app.RequestIDMiddleware,
	)

	api := http.NewServeMux()
	api.Handle(
		"POST /v1/auth/register",
		app.RateLimiterMiddleware(http.HandlerFunc(app.Register)),
	)
	api.Handle(
		"POST /v1/auth/login",
		app.RateLimiterMiddleware(http.HandlerFunc(app.Login)),
	)
	api.Handle(
		"POST /v1/auth/refresh",
		http.HandlerFunc(app.Refresh),
	)
	api.Handle(
		"POST /v1/auth/logout",
		http.HandlerFunc(app.Logout),
	)

	api.Handle(
		"POST /v1/summarize/file",
		app.AuthMiddleware(http.HandlerFunc(app.SummarizeFile)),
	)
	api.Handle(
		"POST /v1/summarize/text",
		app.AuthMiddleware(http.HandlerFunc(app.SummarizeText)),
	)

	api.Handle(
		"GET /v1/tldrs",
		app.AuthMiddleware(http.HandlerFunc(app.UserGetTLDRs)),
	)
	api.Handle(
		"GET /v1/tldrs/{id}",
		app.AuthMiddleware(http.HandlerFunc(app.UserGetTLDR)),
	)
	api.Handle(
		"PATCH /v1/tldrs/{id}",
		app.AuthMiddleware(http.HandlerFunc(app.UserUpdateTLDR)),
	)
	api.Handle(
		"DELETE /v1/tldrs/{id}",
		app.AuthMiddleware(http.HandlerFunc(app.UserDeleteTLDR)),
	)
	api.Handle(
		"DELETE /v1/tldrs",
		app.AuthMiddleware(http.HandlerFunc(app.UserDeleteTLDRs)),
	)

	api.Handle(
		"GET /v1/me",
		app.AuthMiddleware(http.HandlerFunc(app.UserGetMe)),
	)
	api.Handle(
		"PATCH /v1/me/change-username",
		app.AuthMiddleware(http.HandlerFunc(app.UserUpdateUsername)),
	)
	api.Handle(
		"PATCH /v1/me/reset-password",
		app.AuthMiddleware(http.HandlerFunc(app.UserUpdatePassword)),
	)

	admin := http.NewServeMux()
	admin.HandleFunc("GET /tldrs", app.AdminGetTLDRs)
	admin.HandleFunc("GET /tldrs/{id}", app.AdminGetTLDR)
	admin.HandleFunc("PATCH /tldrs/{id}", app.AdminUpdateTLDR)
	admin.HandleFunc("DELETE /tldrs/{id}", app.AdminDeleteTLDR)
	admin.HandleFunc("DELETE /tldrs", app.AdminDeleteTLDRs)
	admin.HandleFunc("DELETE /tldrs/all", app.AdminDeleteAllTLDRs)

	admin.HandleFunc("GET /users", app.AdminGetUsers)
	admin.HandleFunc("GET /users/{id}", app.AdminGetUser)
	admin.HandleFunc("PATCH /users/{id}/change-username", app.AdminUpdateUsername)
	admin.HandleFunc("PATCH /users/{id}/reset-password", app.AdminUpdatePassword)
	admin.HandleFunc("DELETE /users/{id}", app.AdminDeleteUser)
	admin.HandleFunc("DELETE /users/all", app.AdminDeleteAllUsers)

	app.Handle("/", http.FileServer(http.Dir("web/dist")))
	app.Handle(
		"/admin/",
		http.StripPrefix(
			"/admin",
			app.AuthMiddleware(app.AdminMiddleware(admin)),
		),
	)
	app.Handle("/api/", http.StripPrefix("/api", api))
	app.HandleFunc("GET /health", app.Health)

	server := &http.Server{
		Addr:    ":" + app.Config.Port,
		Handler: app.Handler,
	}

	app.Logger.Info(
		"Server started:",
		"port",
		app.Config.Port,
		"environment",
		app.Config.Environment,
	)
	go func() {
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			app.Logger.Error("Error occurred:", "error", err.Error())
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()
	app.Logger.Info("Server shutting down...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		app.Logger.Error("Error occurred shutting down:", "error", err.Error())
	}
}
