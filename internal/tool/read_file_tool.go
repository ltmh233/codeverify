package tool

import (
	"codeverify/internal/core"
	"context"
	"os"
)

type ReadFileTool struct {
	RepoPath string
}

func (t *ReadFileTool) Name() core.ToolName {
	return core.ToolReadFile
}

func (t *ReadFileTool) Execute(ctx context.Context, args map[string]string) (core.Observation, error) {

	path := t.RepoPath + "/" + args["filename"]

	data, err := os.ReadFile(path)
	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, err
	}

	return core.Observation{
		Success: true,
		Output:  string(data),
	}, nil
}
