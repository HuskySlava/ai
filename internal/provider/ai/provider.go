package ai

import "context"

type Provider interface {
	Rewrite(ctx context.Context, text string) (string, error)
	Translate(ctx context.Context, text string, toLanguage string) (string, error)
	General(ctx context.Context, text string) (string, error)
	SendRequest(ctx context.Context, prompt string) (string, error)
}
