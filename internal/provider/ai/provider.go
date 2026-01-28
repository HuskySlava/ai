package ai

import (
	"ai/internal/config"
	"context"
	"fmt"
)

type Provider interface {
	Rewrite(ctx context.Context, text string) (string, error)
	Translate(ctx context.Context, text string, toLanguage string) (string, error)
	Summarize(ctx context.Context, text string) (string, error)
	General(ctx context.Context, text string) (string, error)
}

type baseProvider struct {
	cfg *config.Config
}

func (b *baseProvider) buildPromptRewrite(text string) string {
	return b.cfg.Prompts.Rewrite + " " + text
}

func (b *baseProvider) buildPromptTranslate(text, toLanguage string) string {
	return fmt.Sprintf(b.cfg.Prompts.Translate, toLanguage) + " " + text
}

func (b *baseProvider) buildPromptSummarize(text string) string {
	return b.cfg.Prompts.Summarize + " " + text
}
