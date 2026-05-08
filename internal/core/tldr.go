package core

import (
	"context"

	"github.com/duanechan/tldr/internal/config"
	"google.golang.org/genai"
)

type TLDR struct {
	Config *config.Config
	Client *genai.Client
	Model  *genai.GenerateContentConfig
}

const prompt = `
You are a document summarizer.
Summarize the documents that are provided to you and extract its key points.`

func New() (*TLDR, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	model := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(prompt, genai.RoleUser),
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{})
	if err != nil {
		return nil, err
	}

	return &TLDR{
		Config: cfg,
		Client: client,
		Model:  model,
	}, nil
}
