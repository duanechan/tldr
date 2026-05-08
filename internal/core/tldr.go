package core

import (
	"context"

	"github.com/duanechan/tldr/internal/config"
	"google.golang.org/genai"
)

type TLDR struct {
	Config *config.Config
	Client *genai.Client
}

func New() (*TLDR, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	client, err := genai.NewClient(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &TLDR{
		Config: cfg,
		Client: client,
	}, nil
}
