package agent

import (
	"codeverify/internal/core"
	"codeverify/internal/tool"
	"context"
	"fmt"
)

func Run(ctx context.Context, a Agent, state *core.State, tools map[string]tool.Tool) error {
	const maxSteps = 10

	for step := range maxSteps {
		decision, err := a.Decide(ctx, state)
		if err != nil {
			return err
		}

		fmt.Printf("[step %d]decision: %+v\n", step+1, decision)

		if decision.Type == core.DecisionFinish {
			return nil
		}

		t, ok := tools[string(decision.ToolCall.Name)]
		if !ok {
			return fmt.Errorf("unknown tools: %s", decision.ToolCall.Name)
		}

		observation, _ := t.Execute(ctx, decision.ToolCall.Args)
		fmt.Printf("[step %d]result: %+v\n", step+1, observation)

		state.History = append(state.History, core.Step{
			Decision:    decision,
			Observation: observation,
		})

	}
	return fmt.Errorf("max steps reached")

}
