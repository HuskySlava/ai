package ai

import (
	"ai/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type GeminiProvider struct {
	baseProvider
	apiKey string
	model  string
	client *http.Client
}

func NewGemini(apiKey string, model string, cfg *config.Config) (*GeminiProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("missing GEMINI_API_KEY environment variable")
	}
	// Create a new GeminiProvider and return its address
	return &GeminiProvider{
		baseProvider: baseProvider{cfg: cfg},
		apiKey:       apiKey,
		model:        model,
		client:       &http.Client{},
	}, nil
}

func (p *GeminiProvider) Rewrite(ctx context.Context, input string) (string, error) {
	return p.sendRequest(ctx, p.buildPromptRewrite(input))
}

func (p *GeminiProvider) Translate(ctx context.Context, input string, toLanguage string) (string, error) {
	return p.sendRequest(ctx, p.buildPromptTranslate(input, toLanguage))
}

func (p *GeminiProvider) Summarize(ctx context.Context, input string) (string, error) {
	return p.sendRequest(ctx, p.buildPromptSummarize(input))
}

func (p *GeminiProvider) General(ctx context.Context, text string) (string, error) {
	return p.sendRequest(ctx, text)
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

func (p *GeminiProvider) sendRequest(ctx context.Context, prompt string) (string, error) {
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
