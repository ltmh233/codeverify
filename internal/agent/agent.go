package agent

import (
	"codeverify/internal/core"
	"context"
)

type Agent interface {
	Decide(ctx context.Context, state *core.State) (core.Decision, error)
}

type FakeAgent struct{}

func (f *FakeAgent) Decide(ctx context.Context, state *core.State) (core.Decision, error) {
	switch len(state.History) {
	case 0:
		return core.Decision{
			Type: core.DecisionToolCall,
			ToolCall: &core.ToolCall{
				Name: core.ToolReadFile,
				Args: map[string]string{"filename": "transfer.go"},
			},
		}, nil
	case 1:
		return core.Decision{
			Type: core.DecisionToolCall,
			ToolCall: &core.ToolCall{
				Name: core.ToolWriteTest,
				Args: map[string]string{
					"filename": "transfer_test.go",
					"code": `
package buggy

import "testing"

func TestTransfer(t *testing.T) {
	a := 1
	b := 2
	want := a + b
	got := transfer(a, b)
	if got != want {
		t.Errorf("got %d but want %d", got, want)
	}
}
`},
			},
		}, nil
	case 2:
		return core.Decision{
			Type: core.DecisionToolCall,
			ToolCall: &core.ToolCall{
				Name: core.ToolRunTest,
				Args: map[string]string{},
			},
		}, nil
	default:
		return core.Decision{
			Type: core.DecisionFinish,
			ToolCall: &core.ToolCall{
				Name: "",
				Args: map[string]string{},
			},
		}, nil
	}
}
