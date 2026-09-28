package tool

import (
	"codeverify/internal/core"
	"context"
	"os"
)

type WriteTestTool struct {
	RepoPath string
}

func (t *WriteTestTool) Name() core.ToolName {
	return core.ToolWriteTest
}

func (t *WriteTestTool) Execute(ctx context.Context, args map[string]string) (core.Observation, error) {

	path := t.RepoPath + "/" + args["filename"]

	err := os.WriteFile(path, []byte(args["code"]), 0o644)

	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, err
	}
	return core.Observation{
		Success: true,
		Output:  "",
	}, nil
}
