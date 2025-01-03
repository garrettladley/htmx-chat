package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Service) Routes(r fiber.Router) {
	r.Route("/", func(r fiber.Router) {
		r.Get("/", s.Index)
		r.Get("/chat", s.Chat)
	})
}
