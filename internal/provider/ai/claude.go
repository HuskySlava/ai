package ai

import (
	"ai/internal/config"
	"net/http"
	"time"
)

type ClaudeProvider struct {
	model  string
	client *http.Client
	cfg    *config.Config
}

func NewClaude(model string) (*ClaudeProvider, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &ClaudeProvider{
		model:  model,
		client: &http.Client{Timeout: time.Duration(cfg.HttpTimeoutSeconds) * time.Second},
		cfg:    cfg,
	}, nil
}
