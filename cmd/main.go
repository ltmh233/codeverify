package main

import (
	"codeverify/internal/agent"
	"codeverify/internal/core"
	"codeverify/internal/tool"
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	a := agent.FakeAgent{}

	state := core.State{
		Goal:    "review code and generator tests",
		History: []core.Step{},
	}

	repo := "./testdata/buggy"
	tools := map[string]tool.Tool{
		"read_file": &tool.ReadFileTool{
			RepoPath: repo,
		},
		"write_test": &tool.WriteTestTool{
			RepoPath: repo,
		},
		"run_test": &tool.RunTestTool{
			RepoPath: repo,
		},
	}

	err := agent.Run(ctx, &a, &state, tools)
	if err != nil {
		fmt.Printf("error occur while agent run: %v", err)
	}
}
