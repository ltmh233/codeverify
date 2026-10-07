package main

import (
	"codeverify/internal/agent"
	"codeverify/internal/core"
	"codeverify/internal/tool"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/openai/openai-go/v3/responses"
)

func main() {
	ctx := context.Background()

	a := agent.NewCLient()

	state := core.State{
		Goal:       "review code and generator tests",
		History:    []core.Step{},
		Transcript: []json.RawMessage{},
	}

	systemInputMessage, err := agent.BuildInputMessageJSON(`
You are a code verification agent.
Do not modify production code.
A bug is confirmed only with executable evidence.
`, responses.EasyInputMessageRoleSystem)
	if err != nil {
		log.Fatal(err)
	}
	state.Transcript = append(state.Transcript, systemInputMessage)
	userInputMessage, err := agent.BuildInputMessageJSON(`
Verify the repository at .
`, responses.EasyInputMessageRoleUser)
	if err != nil {
		log.Fatal(err)
	}
	state.Transcript = append(state.Transcript, userInputMessage)

	repoPath := "./testdata/buggy"
	repoAbs, err := filepath.Abs(repoPath)
	if err != nil {
		log.Fatal(err)
	}
	runCtx := core.RunContext{
		RepoPath: repoPath,
		RepoAbs:  repoAbs,
	}
	tools := map[string]tool.Tool{
		"read_file":  &tool.ReadFileTool{},
		"write_file": &tool.WriteTestTool{},
		"run_test":   &tool.RunTestTool{},
		"read_dir":   &tool.ReadDirTool{},
	}
	for _, tool := range tools {
		a.RegisterTool(tool.Spec())
	}

	err = agent.Run(ctx, a, &state, runCtx, tools)
	if err != nil {
		fmt.Printf("error occur while agent run: %v", err)
	}
	fmt.Printf("%v", state.History)
}
