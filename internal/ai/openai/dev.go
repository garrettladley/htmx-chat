//go:build dev

package openai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/garrettladley/htmx-chat/internal/ai"
)

var words = [26]string{
	"apple", "banana", "cherry", "dragon", "elephant",
	"falcon", "grape", "honey", "igloo", "jacket",
	"kite", "lemon", "mango", "ninja", "orange",
	"penguin", "quill", "rabbit", "snake", "tiger",
	"umbrella", "vampire", "whale", "xylophone", "yellow",
	"zebra",
}

func (o *OpenAI) ChatCompletion(ctx context.Context, content string) (<-chan ai.Result, error) {
	results := make(chan ai.Result)
	go func() {
		defer close(results)
		var (
			builder strings.Builder
			ticker  = time.NewTicker(20 * time.Millisecond)
		)
		defer ticker.Stop()
		for i := 0; i < int(o.maxTokens); i++ {
			select {
			case <-ctx.Done():
				results <- ai.Result{Err: ctx.Err()}
				return
			case <-ticker.C:
				word := words[i%len(words)]
				if _, err := builder.WriteString(word + " "); err != nil {
					results <- ai.Result{Err: fmt.Errorf("%w: %v", ErrMessageBuild, err)}
					return
				}
				results <- ai.Result{Message: builder.String()}
			}
		}
	}()
	return results, nil
}
