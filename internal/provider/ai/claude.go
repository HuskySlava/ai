package ai

import (
	"ai/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ClaudeProvider struct {
	apiKey string
	model  string
	client *http.Client
	cfg    *config.Config
}

func NewClaude(apiKey string, model string) (*ClaudeProvider, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &ClaudeProvider{
		apiKey: apiKey,
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

type message struct {
	Role    string
	Content string
}
type claudeRequest struct {
	Model     string
	Messages  []message
	MaxTokens int
}

func (p *ClaudeProvider) SendRequest(ctx context.Context, prompt string) (string, error) {
	url := p.cfg.BaseEndpoints.Claude

	payload := claudeRequest{
		Model: p.model,
		Messages: []message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 1024, // TODO: Configurable
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP Request with Context
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01") // TODO: Configurable
	req.Header.Set("x-api-key", p.apiKey)             // TODO: Configurable

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status: %s", resp.Status)
	}

	return "", nil
}
