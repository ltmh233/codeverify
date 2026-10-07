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

type ReadDirTool struct{}

func (t *ReadDirTool) Name() core.ToolName {
	return core.ToolReadDir
}

func (t *ReadDirTool) Execute(ctx context.Context, runCtx core.RunContext, rawArgs string) (core.Observation, error) {
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

	entrys, err := os.ReadDir(targetAbs)
	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, nil
	}
	var output strings.Builder
	for i, entry := range entrys {
		if i != len(entrys)-1 {
			if entry.IsDir() {
				fmt.Fprintf(&output, "[dir]%s/\n", entry.Name())
			} else {
				fmt.Fprintf(&output, "[file]%s\n", entry.Name())
			}
		} else {
			if entry.IsDir() {
				fmt.Fprintf(&output, "[dir]%s/", entry.Name())
			} else {
				fmt.Fprintf(&output, "[file]%s", entry.Name())
			}
		}
	}

	return core.Observation{
		Success: true,
		Output:  output.String(),
	}, nil
}

func (t *ReadDirTool) Spec() *responses.ToolUnionParam {
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
			Name:        "read_dir",
			Description: openai.String("Read the detials of the dir of current repository"),
		},
	}
}
