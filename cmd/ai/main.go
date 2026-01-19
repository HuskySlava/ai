package main

import (
	"ai/internal/cli"
	"ai/internal/config"
	"ai/internal/provider/ai"
	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const defaultTargetLanguage = "English"

func runModel(model ai.Provider, ctx context.Context, flags *cli.CMDFlags, cfg *config.Config) (string, error) {

	var (
		input string
		err   error
	)

	// Check for file input
	if flags.File != "" {
		input, err = cli.ReadFile(flags.File, cfg.InputFileLimitKB)
		if err != nil {
			return "", err
		}

		// Concatenate text input with file
		if flags.Input != "" {
			input = flags.Input + "\n" + input
		}
	} else {
		input = flags.Input
	}

	if input == "" {
		return "", fmt.Errorf("missing input")
	}

	// Handle conditional flags
	switch {
	case flags.IsRewrite:
		return model.Rewrite(ctx, input)
	case flags.IsTranslate:
		lang := flags.Language
		if lang == "" {
			lang = defaultTargetLanguage
		}
		return model.Translate(ctx, input, lang)
	case flags.IsSummarize:
		return model.Summarize(ctx, input)
	default:
		return model.General(ctx, input)
	}
}

func main() {
	// Load environment variables for (dev only)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using system vars")
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Set CMD flags
	cmdFlags := cli.SetFlags()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.HttpTimeoutSeconds)*time.Second)
	defer cancel()

	models := map[string]func() (ai.Provider, error){
		"gemini": func() (ai.Provider, error) { return ai.NewGemini(os.Getenv("GEMINI_API_KEY"), cfg.Models.Gemini) },
		"openai": func() (ai.Provider, error) { return ai.NewOpenai(os.Getenv("OPENAI_API_KEY"), cfg.Models.Openai) },
		"claude": func() (ai.Provider, error) { return ai.NewClaude(os.Getenv("CLAUDE_API_KEY"), cfg.Models.Claude) },
		"ollama": func() (ai.Provider, error) { return ai.NewOllama(cfg.Models.Ollama) },
		"":       func() (ai.Provider, error) { return ai.NewOllama(cfg.Models.Ollama) },
	}

	newModel, ok := models[cmdFlags.Provider]
	if !ok {
		log.Fatal("Model not implemented")
	}

	model, err := newModel()
	if err != nil {
		log.Fatalf("Error creating model: %v", err)
	}

	res, err := runModel(model, ctx, cmdFlags, cfg)
	if err != nil {
		log.Fatalf("Error running model: %v", err)
		return
	}

	// Copy to clipboard
	if cmdFlags.IsClipboard {
		err = clipboard.WriteAll(res)
		if err != nil {
			fmt.Println("Error copying to clipboard:", err)
		}
	}

	// Output to a file or standard output (stdout)
	if cmdFlags.ToFile != "" {
		err = cli.WriteFile(cmdFlags.ToFile, []byte(res))
		if err != nil {
			log.Fatalf("Error writing file: %v", err)
		}
	} else {
		const cyberCyan = "\033[96m" // bright cyan
		const reset = "\033[0m"

		fmt.Println("\n" + cyberCyan + res + reset + "\n")
	}

}
