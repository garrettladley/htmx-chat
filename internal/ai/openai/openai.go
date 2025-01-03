package openai

import (
	"errors"

	"github.com/garrettladley/htmx-chat/internal/ai"
	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	client    *openai.Client
	maxTokens uint
}

func New(conf ai.Config) *OpenAI {
	return &OpenAI{
		client:    openai.NewClient(conf.APIKey),
		maxTokens: conf.MaxTokens,
	}
}

var (
	ErrStreamCreation = errors.New("failed to create chat completion stream")
	ErrMessageBuild   = errors.New("failed to build message")
)
