package tool

import (
	"codeverify/internal/core"
	"context"
	"encoding/json"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type ReadFileTool struct{}

func (t *ReadFileTool) Name() core.ToolName {
	return core.ToolReadFile
}

func (t *ReadFileTool) Execute(ctx context.Context, runCtx core.RunContext, rawArgs string) (core.Observation, error) {
	var args struct {
		Path string `json:"path"`
	}

	err := json.Unmarshal([]byte(rawArgs), &args)
	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, nil
	}

	if args.Path == "" {
		return core.Observation{
			Success: false,
			Output:  "have not pass argument path",
		}, nil
	}
	targetAbs, err := core.ResolveRepoPath(runCtx.RepoAbs, args.Path)
	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, nil
	}

	data, err := os.ReadFile(targetAbs)
	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, nil
	}

	return core.Observation{
		Success: true,
		Output:  string(data),
	}, nil
}

func (t *ReadFileTool) Spec() *responses.ToolUnionParam {
	return &responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"path": map[string]any{
						"type": "string",
					},
				},
				"required":             []string{"path"},
				"additionalProperties": false,
			},
			Name:        "read_file",
			Description: openai.String("Read a file from the repository by path"),
		},
	}
}
