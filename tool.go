package easyllm

import (
	"context"
)

type Tool interface {
	Name() string

	Description() string

	InputSchema() any

	OutputSchema() any

	Run(ctx context.Context, input any) (any, error)

	Usage() string
}

type ToolCall struct {
	ID           string
	Name         string
	Input        any
	Output       any
	ErrorMessage *string
}
