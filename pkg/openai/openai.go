package openai

import (
	"context"

	"github.com/roadmap-thesis/backend/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
}

func NewClient() *Client {
	client := openai.NewClient(config.OpenAiAPIKey())

	return &Client{
		client: client,
	}
}

type ChatPrompt struct {
	System string
	User   string
}

func (o *Client) Chat(ctx context.Context, prompt ChatPrompt) (*openai.ChatCompletionResponse, error) {
	response, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: config.OpenAiModel(),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt.System,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt.User,
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
