package tool

import (
	"codeverify/internal/core"
	"context"
)

type Tool interface {
	Name() core.ToolName

	Execute(ctx context.Context, args map[string]string) (core.Observation, error)
}
