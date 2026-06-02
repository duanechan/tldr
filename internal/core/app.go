package core

import (
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
	startedAt time.Time
	mu        *sync.RWMutex
	clients   map[string]*rate.Limiter
	*Dependencies
}

type Dependencies struct {
	Config  *config.Config
	DB      *sql.DB
	Queries *database.Queries
	AI      AIModel
	Logger  *slog.Logger
}

const Prompt = `
	You are a document summarizer.

	Summarize the documents, text files, raw text, or images that are provided
	to you.

	Extract its keypoints and ensure the user understands what the contents are.

	Keep the content brief and simple (10% of original content with a maximum
	of 200 words).

	Contents are labelled 'safe', if it's public info, generic documents, code,
	articles.

	Contents are labelled 'mild', if it's personal but not harmful if exposed.

	Contents are labelled 'dangerous' if it contains sensitive information
	(medical records, PII, financial data, credentials, legal docs, etc.)

	Don't make the title too verbose.

	Do not ask for more input.
	`

var ContentConfig = &genai.GenerateContentConfig{
	SystemInstruction: genai.NewContentFromText(Prompt, genai.RoleModel),
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

func Bootstrap() (*Dependencies, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", "./tldr.db?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}

	aiClient, err := NewAIClient(cfg.APIKey)
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

	return &Dependencies{
		Config:  cfg,
		DB:      db,
		Queries: database.New(db),
		AI:      aiClient,
		Logger:  logger,
	}, nil
}

func New(d *Dependencies) *App {

	mux := http.NewServeMux()

	return &App{
		mux:          mux,
		Handler:      mux,
		Dependencies: d,
		startedAt:    time.Now(),
		mu:           &sync.RWMutex{},
		clients:      make(map[string]*rate.Limiter),
	}
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
	a.DB.Close()
}
