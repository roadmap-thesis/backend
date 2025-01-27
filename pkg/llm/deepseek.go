package llm

import (
	"context"
	"fmt"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
	"github.com/roadmap-thesis/backend/pkg/config"
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
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("deepseek: no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}
