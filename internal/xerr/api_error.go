package xerr

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/garrettladley/htmx-chat/internal/xslog"
	"github.com/gofiber/fiber/v2"
)

type APIError struct {
	StatusCode int `json:"statusCode"`
	Message    any `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d %v", e.StatusCode, e.Message)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Message:    err.Error(),
	}
}

func InternalServerError() APIError {
	return NewAPIError(http.StatusInternalServerError, errors.New("internal server error"))
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	var apiErr APIError
	if castedErr, ok := err.(APIError); ok {
		apiErr = castedErr
	} else {
		apiErr = InternalServerError()
	}

	slog.LogAttrs(
		c.Context(),
		slog.LevelError,
		"Error handling request",
		xslog.Error(err),
	)

	return c.Status(apiErr.StatusCode).JSON(apiErr)
}
