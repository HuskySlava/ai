package ai

import (
	"ai/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaProvider struct {
	baseProvider
	model  string
	client *http.Client
}

func NewOllama(model string, cfg *config.Config) (*OllamaProvider, error) {
	return &OllamaProvider{
		baseProvider: baseProvider{cfg: cfg},
		model:        model,
		client:       &http.Client{},
	}, nil
}

func (p *OllamaProvider) Rewrite(ctx context.Context, input string) (string, error) {
	return p.sendRequest(ctx, p.buildPromptRewrite(input))
}

func (p *OllamaProvider) Translate(ctx context.Context, input string, toLanguage string) (string, error) {
	return p.sendRequest(ctx, p.buildPromptTranslate(input, toLanguage))
}

func (p *OllamaProvider) Summarize(ctx context.Context, input string) (string, error) {
	return p.sendRequest(ctx, p.buildPromptSummarize(input))
}

func (p *OllamaProvider) General(ctx context.Context, input string) (string, error) {
	return p.sendRequest(ctx, input)
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"` // This contains the token/word
	Done      bool   `json:"done"`
}

func (p *OllamaProvider) sendRequest(ctx context.Context, prompt string) (string, error) {

	url := p.cfg.BaseEndpoints.Ollama

	payload := ollamaRequest{
		Model:  p.model,
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
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Response, nil
}
