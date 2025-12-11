package main

import (
	"ai/internal/config"
	"ai/internal/provider/ai"
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type CMDFlags struct {
	isRewrite   bool
	isTranslate bool
	model       string
	text        string
}

func setFlags() *CMDFlags {
	flags := &CMDFlags{}

	flag.BoolVar(&flags.isRewrite, "rewrite", false, "AI rewrite function flag")
	flag.BoolVar(&flags.isRewrite, "r", false, "AI rewrite function flag (shorthand)")

	flag.BoolVar(&flags.isTranslate, "translate", false, "AI translate function flag")
	flag.BoolVar(&flags.isTranslate, "t", false, "AI translate function flag (shorthand)")

	flag.StringVar(&flags.model, "model", "", "AI model provider flag")
	flag.StringVar(&flags.model, "m", "", "AI model provider flag (shorthand)")

	flag.StringVar(&flags.text, "prompt", "", "AI prompt")
	flag.StringVar(&flags.text, "p", "", "AI prompt (shorthand)")

	flag.Parse()

	return flags
}

func runModel(model ai.Provider, ctx context.Context, flags *CMDFlags) (string, error) {
	var (
		res string
		err error
	)

	if flags.isRewrite {
		res, err = model.Rewrite(ctx, flags.text)
	} else if flags.isTranslate {
		res, err = model.Translate(ctx, flags.text)
	} else {
		res, err = model.Test(ctx, flags.text)
	}

	return res, err
}

func main() {
	// Load environment variables for (dev only)
	err := godotenv.Load()
	if err != nil {
		log.Println("Error reading .env file")
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Set CMD flags
	cmdFlags := setFlags()

	if cmdFlags.text == "" {
		log.Println("Missing prompt")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.HttpTimeoutSeconds)*time.Second)
	defer cancel()

	var model ai.Provider

	switch cmdFlags.model {
	case "gemini":
		model, err = ai.NewGemini(os.Getenv("GEMINI_API_KEY"), cfg.Models.Gemini)
	case "openai":
		model, err = ai.NewOpenai(os.Getenv("OPENAI_API_KEY"), cfg.Models.Openai)
	case "ollama", "": // Defaults to Ollama if no model flag
		model, err = ai.NewOllama(cfg.Models.Ollama)
	default:
		log.Fatal("Model not implemented")
	}

	if err != nil {
		log.Fatal("Error", err)
	}

	res, err := runModel(model, ctx, cmdFlags)
	if err != nil {
		log.Fatal("Error", err)
		return
	}

	fmt.Println("\n" + res + "\n")

}
