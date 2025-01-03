package server

import (
	"log/slog"

	"github.com/garrettladley/htmx-chat/internal/ai"
	"github.com/garrettladley/htmx-chat/internal/handlers"
	"github.com/garrettladley/htmx-chat/internal/settings"
	"github.com/garrettladley/htmx-chat/internal/xerr"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"

	_ "embed"
)

type Config struct {
	Settings *settings.Settings
	Logger   *slog.Logger
	StaticFn func(*fiber.App)
}

func New(cfg *Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: xerr.ErrorHandler,
	})
	setupMiddleware(app)
	setupHealthCheck(app)
	setupFavicon(app)

	service := handlers.NewService(ai.Config{
		APIKey:    cfg.Settings.APIKey,
		MaxTokens: cfg.Settings.MaxTokens,
	})
	service.Routes(app)
	cfg.StaticFn(app)

	return app
}

func setupMiddleware(app *fiber.App) {
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
}

func setupHealthCheck(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
}

func setupFavicon(app *fiber.App) {
	app.Get("/favicon.ico", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })
}
