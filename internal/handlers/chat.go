package handlers

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) Chat(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")

	results, err := s.ai.ChatCompletion(context.Background(), "Write a poem about why you love Golang")
	if err != nil {
		return err
	}

	c.Context().Response.SetBodyStreamWriter(func(w *bufio.Writer) {
		replacer := strings.NewReplacer("\n", " ", "\r", " ")
		write(w, "data:<div><div>\n\n")
		var prev string
		for {
			select {
			case result, ok := <-results:
				if !ok {
					if err := write(w, none(prev)); err != nil {
						slog.Error("Error writing message", "error", err)
						return
					}
					if err := w.Flush(); err != nil {
						slog.Error("Error flushing writer", "error", err)
						return
					}
					return
				}

				if result.Err != nil {
					if errors.Is(result.Err, context.Canceled) || errors.Is(result.Err, context.DeadlineExceeded) {
						slog.Error("Context error", "error", result.Err)
						return
					}
					slog.Error("Error processing result", "error", result.Err)
					return
				}

				message := replacer.Replace(result.Message)

				prev = result.Message
				if result.Message == "" {
					continue
				}
				if err := write(w, some(message)); err != nil {
					slog.Error("Error writing message", "error", err)
					return
				}

				if err := write(w, "data:<div>\n\n"); err != nil {
					slog.Error("Error writing message", "error", err)
					return
				}

				if err := w.Flush(); err != nil {
					slog.Error("Error flushing writer", "error", err)
					return
				}
			}
		}
	})

	return nil
}

func write(w *bufio.Writer, data string) error {
	_, err := fmt.Fprint(w, data)
	return err
}

func some(content string) string {
	return fmt.Sprintf("data:<div><p>%s</p>\n", content)
}

func none(content string) string {
	return fmt.Sprintf("data:<div id=\"sse-listener\" hx-swap-oob=\"true\"></div>\ndata:<div hx-swap-oob=\"outerHTML:#message-container\"><p>%s</p>\ndata:</div>\n\n", content)
}
