package core

import (
	"github.com/duanechan/tldr/internal/config"
)

type TLDR struct {
	Config config.Config
}

func New() (*TLDR, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	return &TLDR{
		Config: *cfg,
	}, nil
}
