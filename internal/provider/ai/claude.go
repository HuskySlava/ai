package ai

import (
	"ai/internal/config"
	"context"
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

func (p *ClaudeProvider) Rewrite(ctx context.Context, text string) (string, error) {
	return "", nil
}

func (p *ClaudeProvider) Translate(ctx context.Context, text string, toLanguage string) (string, error) {
	return "", nil
}

func (p *ClaudeProvider) Summarize(ctx context.Context, text string) (string, error) {
	return "", nil
}

func (p *ClaudeProvider) General(ctx context.Context, text string) (string, error) {
	return "", nil
}

func (p *ClaudeProvider) SendRequest(ctx context.Context, prompt string) (string, error) {
	return "", nil
}
