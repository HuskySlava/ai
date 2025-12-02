package main

import (
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

	gemini := ai.NewGemini(os.Getenv("GEMINI_API_KEY"), "lol")

	bg := context.Background()
	ctx, cancel := context.WithTimeout(bg, 30*time.Second)

	defer cancel()

	result, err := gemini.Rewrite(ctx, "Hello")
	if err != nil {
		log.Fatal("AI Failer", err)
	}

	log.Println(result)
}
