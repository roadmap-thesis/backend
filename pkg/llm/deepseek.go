package llm

import (
	"context"
	"fmt"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
	"github.com/roadmap-thesis/backend/pkg/config"
	"go.opentelemetry.io/otel/attribute"
)

type DeepSeekClient struct {
	client *deepseek.Client
}

func NewDeepSeekClient() Client {
	client := deepseek.NewClient(config.DeepSeekAPIKey())

	return &DeepSeekClient{
		client: client,
	}
}

func (d *DeepSeekClient) Chat(ctx context.Context, prompt ChatPrompt) (string, error) {
	ctx, span := tracer.Start(ctx, "(*openAiClient.Chat)")
	defer span.End()

	response, err := d.client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model: config.DeepSeekModel(),
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    constants.ChatMessageRoleSystem,
				Content: prompt.System,
			},
			{
				Role:    constants.ChatMessageRoleUser,
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
		return "", fmt.Errorf("deepseek: no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}
