package core

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/duanechan/tldr/internal/config"
	"github.com/duanechan/tldr/internal/database"
	"golang.org/x/time/rate"
	"google.golang.org/genai"
)

const baseURL = "http://localhost:8080"
const apiURL = baseURL + "/api/v1"

type mockAI struct {
	ListFn func(
		ctx context.Context,
		config *genai.ListModelsConfig,
	) (genai.Page[genai.Model], error)
	GenerateContentFn func(
		ctx context.Context,
		model string,
		contents []*genai.Content,
		config *genai.GenerateContentConfig,
	) (*genai.GenerateContentResponse, error)
}

func (m *mockAI) GenerateContent(
	ctx context.Context,
	model string,
	contents []*genai.Content,
	config *genai.GenerateContentConfig,
) (*genai.GenerateContentResponse, error) {
	if m.GenerateContentFn == nil {
		return nil, nil
	}
	return m.GenerateContentFn(ctx, model, contents, config)
}

func (m *mockAI) List(
	ctx context.Context,
	config *genai.ListModelsConfig,
) (genai.Page[genai.Model], error) {
	if m.ListFn == nil {
		return genai.Page[genai.Model]{}, nil
	}
	return m.ListFn(ctx, config)
}

func newTestDependencies(ai AIModel) *Dependencies {
	db, _ := sql.Open("sqlite", ":memory:")
	return &Dependencies{
		Config: &config.Config{
			Port:          "8080",
			Environment:   "test",
			JWTSecret:     "SeCrEtKeY1!2@3#",
			JWTExpiry:     time.Minute * 5,
			RefreshExpiry: time.Hour,
			APIModel:      "foo-1.0",
			APIKey:        "API-Key-1234567890",
			LogLevel:      "info",
		},
		DB:      db,
		Queries: database.New(db),
		AI:      ai,
		Logger:  slog.Default(),
	}
}

func newMockApp() *App {
	return &App{
		mux:          http.NewServeMux(),
		Handler:      http.NewServeMux(),
		startedAt:    time.Now(),
		mu:           &sync.RWMutex{},
		clients:      make(map[string]*rate.Limiter),
		Dependencies: newTestDependencies(&mockAI{}),
	}
}
