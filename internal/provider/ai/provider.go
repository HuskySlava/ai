package ai

import "context"

type Provider interface {
	Rewrite(ctx context.Context, text string) (string, error)
	Translate(ctx context.Context, text string) (string, error)
	sendRequest(ctx context.Context, prompt string) (string, error)
}
