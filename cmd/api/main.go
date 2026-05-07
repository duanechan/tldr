package main

import (
	"log"
	"net/http"

	"github.com/duanechan/tldr/internal/core"
)

func main() {
	app, err := core.New()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", app.Health)

	log.Printf("Server listening on port :%s\n", app.Config.Port)
	log.Printf("Server environment: %s\n", app.Config.Environment)
	if err := http.ListenAndServe(":"+app.Config.Port, mux); err != nil {
		log.Fatal(err)
	}
}
