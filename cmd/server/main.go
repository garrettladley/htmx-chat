package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/garrettladley/htmx-chat/internal/server"
	"github.com/garrettladley/htmx-chat/internal/settings"
	"github.com/garrettladley/htmx-chat/internal/xslog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())

	settings, err := settings.Load()
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to load settings",
			xslog.Error(err),
		)
		os.Exit(1)
	}

	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to connect to database",
			xslog.Error(err),
		)
		os.Exit(1)
	}

	app := server.New(&server.Config{
		Settings: &settings,
		Logger:   logger,
		StaticFn: static,
	})

	go func() {
		if err := app.Listen(":" + settings.App.Port); err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"failed to start server",
				xslog.Error(err),
			)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.LogAttrs(
		ctx,
		slog.LevelInfo,
		"stopping server",
	)
	cancel()

	if err := app.Shutdown(); err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to shutdown server",
			xslog.Error(err),
		)
	}

	slog.LogAttrs(
		ctx,
		slog.LevelInfo,
		"server shutdown",
	)
}

func static(app *fiber.App) {
	app.Get("/public/*", adaptor.HTTPHandler(public()))
	app.Get("/deps/*", adaptor.HTTPHandler(deps()))
}
