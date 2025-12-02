package ai

import "context"

type Provider interface {
	Rewrite(ctx context.Context, text string) (string, error)
}
