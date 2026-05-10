package core

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/duanechan/tldr/internal/config"
	"github.com/lmittmann/tint"
	"google.golang.org/genai"
)

type TLDR struct {
	Config    *config.Config
	Client    *genai.Client
	Model     *genai.GenerateContentConfig
	Logger    *slog.Logger
	startedAt time.Time
}

const prompt = `
You are a document summarizer.
Summarize the documents, text files, raw text (even if brief), or images that are provided to you.
Extract its keypoints, if necessary, and ensure the user understands what the contents are.
Do not ask for more input.`

func New() (*TLDR, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	model := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(prompt, genai.RoleModel),
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{})
	if err != nil {
		return nil, err
	}

	var logLevel slog.Leveler
	switch cfg.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{Level: logLevel}))

	return &TLDR{
		Config:    cfg,
		Client:    client,
		Model:     model,
		Logger:    logger,
		startedAt: time.Now(),
	}, nil
}
