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

type OllamaProvider struct {
	model  string
	client *http.Client
	cfg    *config.Config
}

func NewOllama(model string) (*OllamaProvider, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &OllamaProvider{
		model:  model,
		client: &http.Client{Timeout: time.Duration(cfg.HttpTimeoutSeconds) * time.Second},
		cfg:    cfg,
	}, nil
}

func (p *OllamaProvider) Rewrite(ctx context.Context, text string) (string, error) {
	prompt := p.cfg.Prompts.Rewrite + " " + text
	return p.SendRequest(ctx, prompt)
}

func (p *OllamaProvider) Translate(ctx context.Context, text string, toLanguage string) (string, error) {
	prompt := fmt.Sprintf(p.cfg.Prompts.Translate, toLanguage) // Inject target language to config prompt
	prompt += " " + text                                       // Add text to be translated
	return p.SendRequest(ctx, prompt)
}

func (p *OllamaProvider) Test(ctx context.Context, text string) (string, error) {
	prompt := text
	return p.SendRequest(ctx, prompt)
}

type ollamaRequest struct {
	Model  string
	Prompt string
	Stream bool
}

type ollamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"` // This contains the token/word
	Done      bool   `json:"done"`
}

func (p *OllamaProvider) SendRequest(ctx context.Context, prompt string) (string, error) {

	url := p.cfg.BaseEndpoints.Ollama

	payload := ollamaRequest{
		Model:  p.cfg.Models.Ollama,
		Prompt: prompt,
		Stream: false,
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

	// Execute
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status: %s", resp.Status)
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Response, err
}
