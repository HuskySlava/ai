package ai

import (
	"ai/internal/config"
	"context"
	"log"
	"net/http"
	"time"
)

type GeminiProvider struct {
	apiKey string
	model  string
	client *http.Client
}

func NewGemini(apiKey string, model string) *GeminiProvider {
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		log.Fatal("Error reading config", err)
	}
	log.Println("Config loaded", *cfg)
	// Create a new GeminiProvider and return its address
	return &GeminiProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *GeminiProvider) Rewrite(ctx context.Context, text string) (string, error) {
	log.Println("Rewriting text...")
	return "[GEMINI RESULT]", nil
}
