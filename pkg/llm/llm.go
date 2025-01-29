package llm

import (
	"context"

	"go.opentelemetry.io/otel"
)

type Client interface {
	Chat(ctx context.Context, prompt ChatPrompt) (string, error)
}

type Provider string

const (
	OpenAI   Provider = "openai"
	DeepSeek Provider = "deepseek"
)

var (
	tracer = otel.Tracer("LLM")
)
