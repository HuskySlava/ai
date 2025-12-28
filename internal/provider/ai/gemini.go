package ai

import (
	"ai/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type GeminiProvider struct {
	apiKey string
	model  string
	client *http.Client
	cfg    *config.Config
}

func NewGemini(apiKey string, model string) (*GeminiProvider, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	// Create a new GeminiProvider and return its address
	return &GeminiProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: time.Duration(cfg.HttpTimeoutSeconds) * time.Second},
		cfg:    cfg,
	}, nil
}

func (p *GeminiProvider) Rewrite(ctx context.Context, text string) (string, error) {
	prompt := p.cfg.Prompts.Rewrite + " " + text
	return p.SendRequest(ctx, prompt)
}

func (p *GeminiProvider) Translate(ctx context.Context, text string, toLanguage string) (string, error) {
	prompt := fmt.Sprintf(p.cfg.Prompts.Translate, toLanguage) // Inject target language to config prompt
	prompt += " " + text                                       // Add text to be translated
	return p.SendRequest(ctx, prompt)
}

func (p *GeminiProvider) General(ctx context.Context, text string) (string, error) {
	prompt := text
	return p.SendRequest(ctx, prompt)
}

type geminiRequest struct {
	Contents []content `json:"contents"`
}

type content struct {
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

// responseBody matches Gemini's response structure
type geminiResponse struct {
	Candidates []struct {
		Content content `json:"content"`
	} `json:"candidates"`
}

func (p *GeminiProvider) SendRequest(ctx context.Context, prompt string) (string, error) {
	// Endpoint construction
	// Example: https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=XYZ
	url := fmt.Sprintf(
		p.cfg.BaseEndpoints.Gemini+"%s:generateContent?key=%s",
		p.model,
		p.apiKey,
	)

	// Prepare JSON payload
	payload := geminiRequest{
		Contents: []content{
			{
				Parts: []part{
					{Text: prompt},
				},
			},
		},
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

	// Parse Response
	var result geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract text safely
	if len(result.Candidates) == 0 ||
		len(result.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no response content from Gemini")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
