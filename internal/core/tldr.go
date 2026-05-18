package core

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/duanechan/tldr/internal/config"
	"github.com/duanechan/tldr/internal/database"
	"github.com/lmittmann/tint"
	"golang.org/x/time/rate"
	"google.golang.org/genai"
)

type App struct {
	mux       *http.ServeMux
	Handler   http.Handler
	db        *sql.DB
	Queries   *database.Queries
	Config    *config.Config
	Client    *genai.Client
	Model     *genai.GenerateContentConfig
	Logger    *slog.Logger
	startedAt time.Time
	mu        *sync.RWMutex
	clients   map[string]*rate.Limiter
}

const prompt = `
	You are a document summarizer.
	
	Summarize the documents, text files, raw text, or images that are provided
	to you.
	
	Extract its keypoints, if necessary, and ensure the user understands what
	the contents are.
	
	Keep the content brief and simple (10% of original content with a maximum
	of 200 words).

	Provide a flag (Safe, Mild, Dangerous) based on the contents of the
	summary. Sensitive information (medical info, credit cards, etc.) should
	be automatically flagged as dangerous.
	
	Don't make the title too verbose.
	
	Do not ask for more input.
	
	No markdown.
	
	ONLY RETURN a response in this JSON format:
	{
		"title": string,
		"content": string,
		"flag": Safe | Mild | Dangerous
	}
	`

func New() (*App, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", "./tldr.db?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}

	model := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(prompt, genai.RoleModel),
		ResponseMIMEType:  "application/json",
		ResponseJsonSchema: map[string]any{
			"type":     "object",
			"required": []string{"title", "content", "flag"},
			"properties": map[string]any{
				"title": map[string]any{
					"type":        "string",
					"description": "The title of the TLDR.",
				},
				"content": map[string]any{
					"type":        "string",
					"description": "The summarized content.",
				},
				"flag": map[string]any{
					"type":        "string",
					"enum":        []string{"Safe", "Mild", "Dangerous"},
					"description": "Safety level of the content.",
				},
			},
		},
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: cfg.APIKey,
	})
	if err != nil {
		return nil, err
	}

	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		return nil, err
	}

	logger := slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{Level: logLevel}),
	)
	mux := http.NewServeMux()

	return &App{
		mux:       mux,
		Handler:   mux,
		db:        db,
		Queries:   database.New(db),
		Config:    cfg,
		Client:    client,
		Model:     model,
		Logger:    logger,
		startedAt: time.Now(),
		mu:        &sync.RWMutex{},
		clients:   make(map[string]*rate.Limiter),
	}, nil
}

func (a *App) Handle(pattern string, handler http.Handler) {
	a.mux.Handle(pattern, handler)
}

func (a *App) HandleFunc(pattern string, handler http.HandlerFunc) {
	a.mux.HandleFunc(pattern, handler)
}

func (a *App) Use(middleware ...func(http.Handler) http.Handler) {
	for _, m := range middleware {
		a.Handler = m(a.Handler)
	}
}

func (a *App) CloseDB() {
	a.Logger.Info("Database connection closed")
	a.db.Close()
}
