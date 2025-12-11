package ai

import (
	"ai/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenaiProvider struct {
	apiKey string
	model  string
	client *http.Client
	cfg    *config.Config
}

func NewOpenai(apiKey string, model string) (*OpenaiProvider, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &OpenaiProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: time.Duration(cfg.HttpTimeoutSeconds) * time.Second},
		cfg:    cfg,
	}, nil
}

func (p *OpenaiProvider) Rewrite(ctx context.Context, text string) (string, error) {
	prompt := p.cfg.Prompts.Rewrite + " " + text
	return p.SendRequest(ctx, prompt)
}

func (p *OpenaiProvider) Translate(ctx context.Context, text string) (string, error) {
	prompt := p.cfg.Prompts.Translate + " " + text
	return p.SendRequest(ctx, prompt)
}

func (p *OpenaiProvider) Test(ctx context.Context, text string) (string, error) {
	return p.SendRequest(ctx, text)
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type ChatChoice struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}

type ChatResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Choices []ChatChoice `json:"choices"`
}

func (p *OpenaiProvider) SendRequest(ctx context.Context, prompt string) (string, error) {

	payload := ChatRequest{
		Model: p.model,
		Messages: []Message{
			{Role: "system", Content: "You are a concise assistant."},
			{Role: "user", Content: prompt},
		},
		Temperature: 1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.cfg.BaseEndpoints.Openai, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("API returned empty choices")
	}

	return result.Choices[0].Message.Content, nil
}
