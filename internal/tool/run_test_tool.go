package tool

import (
	"codeverify/internal/core"
	"context"
	"os/exec"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type RunTestTool struct{}

func (t *RunTestTool) Name() core.ToolName {
	return core.ToolRunTest
}

func (t *RunTestTool) Execute(ctx context.Context, runCtx core.RunContext, rawArgs string) (core.Observation, error) {

	cmd := exec.CommandContext(ctx, "go", "test", "./...")
	cmd.Dir = runCtx.RepoAbs

	output, err := cmd.CombinedOutput()

	return core.Observation{
		Success: err == nil,
		Output:  string(output),
	}, nil
}

func (t *RunTestTool) Spec() *responses.ToolUnionParam {
	return &responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Parameters: map[string]any{
				"type":                 "object",
				"properties":           map[string]any{},
				"additionalProperties": false,
			},
			Name:        string(core.ToolRunTest),
			Description: openai.String("Run all go test in current repository"),
		},
	}
}
