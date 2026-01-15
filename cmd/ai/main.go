package main

import (
	"ai/internal/config"
	"ai/internal/provider/ai"
	"context"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const defaultTargetLanguage = "English"

type CMDFlags struct {
	isRewrite   bool
	isTranslate bool
	isSummarize bool
	isClipboard bool
	provider    string
	input       string
	language    string
}

func setFlags() *CMDFlags {
	flags := &CMDFlags{}

	var rewrite, r bool
	var translate, t bool
	var summarize, s bool
	var copyClipboard, c bool
	var provider, p string
	var input, i string
	var language, l string

	flag.BoolVar(&rewrite, "rewrite", false, "AI rewrite function flag")
	flag.BoolVar(&r, "r", false, "AI rewrite function flag (shorthand)")

	flag.BoolVar(&translate, "translate", false, "AI translate function flag")
	flag.BoolVar(&t, "t", false, "AI translate function flag (shorthand)")

	flag.BoolVar(&summarize, "summarize", false, "AI summarize function flag")
	flag.BoolVar(&s, "s", false, "AI summarize function flag (shorthand)")

	flag.BoolVar(&copyClipboard, "clipboard", false, "Copy result to clipboard automatically")
	flag.BoolVar(&c, "c", false, "Copy result to clipboard automatically (shorthand)")

	flag.StringVar(&provider, "provider", "", "AI model provider flag")
	flag.StringVar(&p, "p", "", "AI model provider flag (shorthand)")

	flag.StringVar(&input, "input", "", "AI prompt")
	flag.StringVar(&i, "i", "", "AI prompt (shorthand)")

	flag.StringVar(&language, "language", "", "Translation target language")
	flag.StringVar(&l, "l", "", "Translation target language (shorthand)")

	flag.Parse()

	firstNonEmpty := func(a, b string) string {
		if a != "" {
			return a
		}
		return b
	}

	flags.isRewrite = rewrite || r
	flags.isTranslate = translate || t
	flags.isSummarize = summarize || s
	flags.isClipboard = copyClipboard || c
	flags.provider = firstNonEmpty(provider, p)
	flags.input = firstNonEmpty(input, i)
	flags.language = firstNonEmpty(language, l)

	return flags
}

func runModel(model ai.Provider, ctx context.Context, flags *CMDFlags) (string, error) {
	switch {
	case flags.isRewrite:
		return model.Rewrite(ctx, flags.input)
	case flags.isTranslate:
		lang := flags.language
		if lang == "" {
			lang = defaultTargetLanguage
		}
		return model.Translate(ctx, flags.input, lang)
	case flags.isSummarize:
		return model.Summarize(ctx, flags.input)
	default:
		return model.General(ctx, flags.input)
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
	cmdFlags := setFlags()

	if cmdFlags.input == "" {
		log.Println("Missing prompt")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.HttpTimeoutSeconds)*time.Second)
	defer cancel()

	models := map[string]func() (ai.Provider, error){
		"gemini": func() (ai.Provider, error) { return ai.NewGemini(os.Getenv("GEMINI_API_KEY"), cfg.Models.Gemini) },
		"openai": func() (ai.Provider, error) { return ai.NewOpenai(os.Getenv("OPENAI_API_KEY"), cfg.Models.Openai) },
		"claude": func() (ai.Provider, error) { return ai.NewClaude(os.Getenv("CLAUDE_API_KEY"), cfg.Models.Claude) },
		"ollama": func() (ai.Provider, error) { return ai.NewOllama(cfg.Models.Ollama) },
		"":       func() (ai.Provider, error) { return ai.NewOllama(cfg.Models.Ollama) },
	}

	newModel, ok := models[cmdFlags.provider]
	if !ok {
		log.Fatal("Model not implemented")
	}

	model, err := newModel()
	if err != nil {
		log.Fatalf("Error creating model: %v", err)
	}

	res, err := runModel(model, ctx, cmdFlags)
	if err != nil {
		log.Fatalf("Error running model: %v", err)
		return
	}

	// Copy to clipboard
	if cmdFlags.isClipboard {
		err = clipboard.WriteAll(res)
		if err != nil {
			fmt.Println("Error copying to clipboard:", err)
		}
	}

	fmt.Println("\n" + res + "\n")

}

// TEST COMMIT
