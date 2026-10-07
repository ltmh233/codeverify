package agent

import (
	"codeverify/internal/core"
	"context"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type Agent interface {
	Decide(ctx context.Context, state *core.State) (core.Turn, error)
	RegisterTool(tool *responses.ToolUnionParam)
}

type LLMAgent struct {
	client openai.Client
	tools  []responses.ToolUnionParam
}

func (a *LLMAgent) Decide(ctx context.Context, state *core.State) (core.Turn, error) {
	res, err := a.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: "deepseek-flash",
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: buildInput(state),
		},
		Tools: a.tools,
	})
	if err != nil {
		return core.Turn{}, err
	}
	return parseResponse(res, state)
}

func (a *LLMAgent) RegisterTool(tool *responses.ToolUnionParam) {
	a.tools = append(a.tools, *tool)
}

func NewCLient() *LLMAgent {
	return &LLMAgent{
		client: openai.NewClient(option.WithAPIKey(os.Getenv("DEEPSEEK_API_KEY")), option.WithBaseURL("https://api.deepseek.com")),
		tools:  []responses.ToolUnionParam{},
	}
}
