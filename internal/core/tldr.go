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

type TLDR struct {
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
	
	Keep the content brief and simple (250-500 words).
	
	Don't make the title too verbose.
	
	Do not ask for more input.
	
	No markdown.
	
	ONLY RETURN a response in this JSON format:
	{
		"title": string,
		"content": string
	}
	`

func New() (*TLDR, error) {
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
			"required": []string{"title", "content"},
			"properties": map[string]any{
				"title": map[string]any{
					"type":        "string",
					"description": "The title of the TLDR.",
				},
				"content": map[string]any{
					"type":        "string",
					"description": "The summarized content.",
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

	logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{Level: logLevel}))
	mux := http.NewServeMux()

	return &TLDR{
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

func (t *TLDR) Handle(pattern string, handler http.Handler) {
	t.mux.Handle(pattern, handler)
}

func (t *TLDR) HandleFunc(pattern string, handler http.HandlerFunc) {
	t.mux.HandleFunc(pattern, handler)
}

func (t *TLDR) Use(middleware ...func(http.Handler) http.Handler) {
	for _, m := range middleware {
		t.Handler = m(t.Handler)
	}
}

func (t *TLDR) CloseDB() {
	t.Logger.Info("Database connection closed")
	t.db.Close()
}
