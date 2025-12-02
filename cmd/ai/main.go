package main

import (
	"ai/internal/config"
	"ai/internal/provider/ai"
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error reading environment variables")
	}

	// Load config
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		log.Fatal("Error loading config")
	}

	gemini := ai.NewGemini(os.Getenv("GEMINI_API_KEY"), "lol")

	bg := context.Background()
	ctx, cancel := context.WithTimeout(bg, time.Duration(cfg.HttpTimeoutSeconds)*time.Second)

	defer cancel()

	result, err := gemini.Rewrite(ctx, "Hello")
	if err != nil {
		log.Fatal("AI Failer", err)
	}

	log.Println(result)
}
