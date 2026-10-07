package tool

import (
	"codeverify/internal/core"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type WriteTestTool struct{}

func (t *WriteTestTool) Name() core.ToolName {
	return core.ToolWriteFile
}

func (t *WriteTestTool) Execute(ctx context.Context, runCtx core.RunContext, rawArgs string) (core.Observation, error) {
	var args struct {
		Path    string `json:"path"`
		Content string `json:"content"`
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
	if !strings.HasSuffix(args.Path, "_test.go") {
		return core.Observation{
			Success: false,
			Output:  "path must end with _test.go",
		}, nil
	}

	path, err := core.ResolveRepoPath(runCtx.RepoAbs, args.Path)
	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, nil
	}

	if args.Content == "" {
		return core.Observation{
			Success: false,
			Output:  "have not pass argument content",
		}, nil

	}
	err = os.WriteFile(path, []byte(args.Content), 0o644)

	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, nil
	}
	return core.Observation{
		Success: true,
		Output:  fmt.Sprintf("write %s success", path),
	}, nil
}

func (t *WriteTestTool) Spec() *responses.ToolUnionParam {
	return &responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"path": map[string]any{
						"type": "string",
					},
					"content": map[string]any{
						"type": "string",
					},
				},
				"required":             []string{"path", "content"},
				"additionalProperties": false,
			},
			Name:        "write_file",
			Description: openai.String("Write a test file by path, path must end with _test.go"),
		},
	}
}
