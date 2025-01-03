//go:build !dev

package openai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/garrettladley/htmx-chat/internal/ai"
	"github.com/sashabaranov/go-openai"
)

func (o *OpenAI) ChatCompletion(ctx context.Context, content string) (<-chan ai.Result, error) {
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: int(o.maxTokens),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful AI assistant.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
		Stream: true,
	}

	stream, err := o.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStreamCreation, err)
	}

	results := make(chan ai.Result)

	go func() {
		defer close(results)
		defer stream.Close()

		var builder strings.Builder
		for {
			select {
			case <-ctx.Done():
				results <- ai.Result{Err: ctx.Err()}
				return
			default:
				response, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					return
				}
				if err != nil {
					results <- ai.Result{Err: fmt.Errorf("stream receive error: %w", err)}
					return
				}

				message := response.Choices[0].Delta.Content
				if _, err := builder.WriteString(message); err != nil {
					results <- ai.Result{Err: fmt.Errorf("%w: %v", ErrMessageBuild, err)}
					return
				}

				results <- ai.Result{Message: builder.String()}
			}
		}
	}()

	return results, nil
}
