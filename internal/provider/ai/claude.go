package ai

import (
	"ai/internal/config"
	"context"
	"fmt"
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
	prompt := p.cfg.Prompts.Rewrite + " " + text
	return p.SendRequest(ctx, prompt)
}

func (p *ClaudeProvider) Translate(ctx context.Context, text string, toLanguage string) (string, error) {
	prompt := fmt.Sprintf(p.cfg.Prompts.Translate, toLanguage) // Inject target language to config prompt
	prompt += " " + text                                       // Add text to be translated
	return p.SendRequest(ctx, prompt)
}

func (p *ClaudeProvider) Summarize(ctx context.Context, text string) (string, error) {
	prompt := fmt.Sprintf(p.cfg.Prompts.Summarize)
	prompt += " " + text
	return p.SendRequest(ctx, prompt)
}

func (p *ClaudeProvider) General(ctx context.Context, text string) (string, error) {
	prompt := text
	return p.SendRequest(ctx, prompt)
}

func (p *ClaudeProvider) SendRequest(ctx context.Context, prompt string) (string, error) {
	return "", nil
}
