package llm

import (
	"context"
	"fmt"

	"github.com/roadmap-thesis/backend/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

type openAiClient struct {
	client *openai.Client
}

func NewOpenAiClient() Client {
	client := openai.NewClient(config.OpenAiAPIKey())

	return &openAiClient{
		client: client,
	}
}

func (o *openAiClient) Chat(ctx context.Context, prompt ChatPrompt) (string, error) {
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
		return "", err
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

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("openai: no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}
