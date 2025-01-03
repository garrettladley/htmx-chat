package handlers

import (
	"github.com/garrettladley/htmx-chat/internal/views/chat"
	"github.com/garrettladley/htmx-chat/internal/xtempl"
	"github.com/gofiber/fiber/v2"
)

func (s *Service) Index(c *fiber.Ctx) error {
	return xtempl.Render(c, chat.Index())
}
