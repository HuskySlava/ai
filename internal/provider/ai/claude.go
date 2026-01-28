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

type ClaudeProvider struct {
	baseProvider
	apiKey string
	model  string
	client *http.Client
}

func NewClaude(apiKey string, model string, cfg *config.Config) (*ClaudeProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("missing CLAUDE_API_KEY environment variable")
	}
	return &ClaudeProvider{
		baseProvider: baseProvider{cfg: cfg},
		apiKey:       apiKey,
		model:        model,
		client:       &http.Client{},
	}, nil
}

func (p *ClaudeProvider) Rewrite(ctx context.Context, input string) (string, error) {
	return p.sendRequest(ctx, p.buildPromptRewrite(input))
}

func (p *ClaudeProvider) Translate(ctx context.Context, input string, toLanguage string) (string, error) {
	return p.sendRequest(ctx, p.buildPromptTranslate(input, toLanguage))
}

func (p *ClaudeProvider) Summarize(ctx context.Context, input string) (string, error) {
	return p.sendRequest(ctx, p.buildPromptSummarize(input))
}

func (p *ClaudeProvider) General(ctx context.Context, input string) (string, error) {
	return p.sendRequest(ctx, input)
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type claudeRequest struct {
	Model     string    `json:"model"`
	Messages  []message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type claudeResponse struct {
	ID           string           `json:"id"`
	Type         string           `json:"type"`
	Role         string           `json:"role"`
	Model        string           `json:"model"`
	Content      []messageContent `json:"content"`
	StopReason   string           `json:"stop_reason"`
	StopSequence *string          `json:"stop_sequence"`
	Usage        usage            `json:"usage"`
}

// Content represents a content block in the response
type messageContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func (p *ClaudeProvider) sendRequest(ctx context.Context, prompt string) (string, error) {
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
	req.Header.Set("x-api-key", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result claudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return result.Content[0].Text, nil
}
