package llm

import "context"

type Client interface {
	Chat(ctx context.Context, prompt ChatPrompt) (string, error)
}
