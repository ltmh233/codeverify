package tool

import (
	"codeverify/internal/core"
	"context"
	"os/exec"
)

type RunTestTool struct {
	RepoPath string
}

func (t *RunTestTool) Name() core.ToolName {
	return core.ToolRunTest
}

func (t *RunTestTool) Execute(ctx context.Context, args map[string]string) (core.Observation, error) {

	cmd := exec.CommandContext(ctx, "go", "test", "./...")
	cmd.Dir = t.RepoPath

	output, err := cmd.CombinedOutput()

	return core.Observation{
		Success: err == nil,
		Output:  string(output),
	}, err
}
