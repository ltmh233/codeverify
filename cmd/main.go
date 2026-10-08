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
你是 CodeVerify, 一个自主代码验证 Agent。

你的任务是分析代码、识别潜在缺陷，并通过执行测试获取证据。

规则：
1. 禁止修改生产代码。
2. 只有获得可靠的可执行证据，才能确认缺陷。
3. 工具调用必须严格遵守对应的 JSON Schema。
4. 所有面向用户的自然语言内容必须使用简体中文。
5. 风险描述、测试分析、验证结论及最终报告必须使用简体中文。
6. 代码、工具名称、文件路径、JSON 字段及原始错误输出保持原样。
7. 当验证任务结束时，使用简体中文生成最终报告。

`, responses.EasyInputMessageRoleSystem)
	if err != nil {
		log.Fatal(err)
	}
	state.Transcript = append(state.Transcript, systemInputMessage)
	userInputMessage, err := agent.BuildInputMessageJSON(`
检验 . 目录下的仓库
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
	fmt.Printf("%v", state.History[len(state.History)-1].Turn.FinalAnswer)
}
