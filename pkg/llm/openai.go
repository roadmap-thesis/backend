package llm

import (
	"context"
	"fmt"

	"github.com/roadmap-thesis/backend/pkg/config"
	"github.com/sashabaranov/go-openai"
	"go.opentelemetry.io/otel/attribute"
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
	ctx, span := tracer.Start(ctx, "(*openAiClient.Chat)")
	defer span.End()

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
		span.RecordError(err)
		return "", err
	}

	span.SetAttributes(
		attribute.String("id", response.ID),
		attribute.String("model", response.Model),
		attribute.String("object", response.Object),
		attribute.Int64("created", response.Created),
		attribute.Int("completion_tokens", response.Usage.CompletionTokens),
		attribute.Int("prompt_tokens", response.Usage.PromptTokens),
		attribute.Int("total_tokens", response.Usage.TotalTokens),
	)

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("openai: no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}
