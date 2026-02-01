package main

import (
	"ai/internal/cli"
	"ai/internal/config"
	"ai/internal/provider/ai"
	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const defaultTargetLanguage = "English"

func readStdin(sizeLimitKB int) (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	// ModeCharDevice is set when stdin is a terminal (not a pipe)
	if info.Mode()&os.ModeCharDevice != 0 {
		return "", nil
	}
	sizeLimit := int64(sizeLimitKB * 1024)
	data, err := io.ReadAll(io.LimitReader(os.Stdin, sizeLimit+1))
	if err != nil {
		return "", fmt.Errorf("failed to read stdin: %w", err)
	}
	if int64(len(data)) > sizeLimit {
		return "", fmt.Errorf("stdin too large (limit %dKB)", sizeLimitKB)
	}
	return string(data), nil
}

func runModel(model ai.Provider, ctx context.Context, flags *cli.CMDFlags, cfg *config.Config) (string, error) {

	var inputParts []string

	if flags.Input != "" {
		inputParts = append(inputParts, flags.Input)
	}

	if flags.File != "" {
		fileContent, err := cli.ReadFile(flags.File, cfg.InputFileLimitKB)
		if err != nil {
			return "", err
		}
		inputParts = append(inputParts, fileContent)
	}

	if len(inputParts) == 0 {
		stdinContent, err := readStdin(cfg.InputFileLimitKB)
		if err != nil {
			return "", err
		}
		if stdinContent != "" {
			inputParts = append(inputParts, stdinContent)
		}
	}

	input := strings.Join(inputParts, "\n")
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
		"gemini": func() (ai.Provider, error) { return ai.NewGemini(os.Getenv("GEMINI_API_KEY"), cfg.Models.Gemini, cfg) },
		"openai": func() (ai.Provider, error) { return ai.NewOpenai(os.Getenv("OPENAI_API_KEY"), cfg.Models.Openai, cfg) },
		"claude": func() (ai.Provider, error) { return ai.NewClaude(os.Getenv("CLAUDE_API_KEY"), cfg.Models.Claude, cfg) },
		"ollama": func() (ai.Provider, error) { return ai.NewOllama(cfg.Models.Ollama, cfg) },
		"":       func() (ai.Provider, error) { return ai.NewOllama(cfg.Models.Ollama, cfg) },
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
