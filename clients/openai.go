package clients

import (
	"context"

	"github.com/HotPotatoC/roadmap_gen/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	client *openai.Client
}

func NewOpenAIClient() *OpenAI {
	client := openai.NewClient(config.OpenAiAPIKey())

	return &OpenAI{
		client: client,
	}
}

func (o *OpenAI) Chat(ctx context.Context, prompt string) (*openai.ChatCompletionResponse, error) {
	response, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: config.OpenAiModel(),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant/teacher/mentor with broad knowledge on many topics.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	log.Info().Dict("openai_request", zerolog.Dict().
		Str("id", response.ID).
		Str("model", response.Model).
		Str("object", response.Object).
		Int64("created", response.Created).
		Dict("usage", zerolog.Dict().
			Int("completion_tokens", response.Usage.CompletionTokens).
			Int("prompt_tokens", response.Usage.PromptTokens).
			Int("total_tokens", response.Usage.TotalTokens),
		),
	).Msg("OpenAI chat request")

	return &response, nil
}
