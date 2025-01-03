package handlers

import (
	"github.com/garrettladley/htmx-chat/internal/ai"
	"github.com/garrettladley/htmx-chat/internal/ai/openai"
)

type Service struct {
	ai ai.Service
}

func NewService(conf ai.Config) *Service {
	return &Service{
		ai: openai.New(conf),
	}
}
