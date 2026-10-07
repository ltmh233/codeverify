package tool

import (
	"codeverify/internal/core"
	"context"

	"github.com/openai/openai-go/v3/responses"
)

type Tool interface {
	Name() core.ToolName

	Execute(ctx context.Context, runCtx core.RunContext, rawArgs string) (core.Observation, error)

	Spec() *responses.ToolUnionParam
}
