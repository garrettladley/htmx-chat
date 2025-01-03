package ai

import "context"

type Service interface {
	ChatCompletion(ctx context.Context, content string) (<-chan Result, error)
}

type Config struct {
	APIKey    string
	MaxTokens uint
}

type Result struct {
	Message string
	Err     error
}
