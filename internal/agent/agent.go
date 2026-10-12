package agent

import (
	"codeverify/internal/core"
	"codeverify/internal/tool"
	"context"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type Agent interface {
	Decide(ctx context.Context, state *core.State, tools tool.Registry) (core.Turn, error)
}

type LLMAgent struct {
	client openai.Client
}

func (a *LLMAgent) Decide(ctx context.Context, state *core.State, tools tool.Registry) (core.Turn, error) {
	res, err := a.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: "deepseek-flash",
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: buildInput(state),
		},
		Tools: tools.Specs(),
	})
	if err != nil {
		return core.Turn{}, err
	}
	return parseResponse(res, state)
}

func NewCLient() *LLMAgent {
	return &LLMAgent{
		client: openai.NewClient(option.WithAPIKey(os.Getenv("DEEPSEEK_API_KEY")), option.WithBaseURL("https://api.deepseek.com")),
	}
}
