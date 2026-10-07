package agent

import (
	"codeverify/internal/core"
	"codeverify/internal/tool"
	"context"
	"fmt"
)

func Run(ctx context.Context, a Agent, state *core.State, runCtx core.RunContext, tools map[string]tool.Tool) error {
	const maxSteps = 100

	for s := range maxSteps {
		turn, err := a.Decide(ctx, state)
		if err != nil {
			return err
		}
		step := core.Step{
			Turn: core.Turn{
				Reasonings: turn.Reasonings,
				ToolCalls:  []*core.ToolCall{},
			},
			ToolResults: []core.ToolResult{},
		}
		for _, toolCall := range turn.ToolCalls {
			fmt.Printf("[step %d]toolCall: %+v\n", s+1, toolCall.Name)
			step.Turn.ToolCalls = append(step.Turn.ToolCalls, toolCall)

			t, ok := tools[string(toolCall.Name)]
			if !ok {
				observation := core.Observation{
					Success: false,
					Output:  fmt.Sprintf("function_call failed: unknown name of function : %s", toolCall.Name),
				}
				step.ToolResults = append(step.ToolResults, core.ToolResult{
					CallID:      toolCall.CallID,
					Observation: observation,
				})
				functionCallOutputJSON, err := buildFunctionCallOutputJSON(toolCall.CallID, observation)
				if err != nil {
					return err
				}

				state.Transcript = append(state.Transcript, functionCallOutputJSON)
			} else {

				observation, err := t.Execute(ctx, runCtx, toolCall.RawArgs)
				if err != nil {
					return err
				}
				observation.Output = core.SanitizePath(observation.Output, runCtx.RepoAbs)
				step.ToolResults = append(step.ToolResults, core.ToolResult{
					CallID:      toolCall.CallID,
					Observation: observation,
				})
				fmt.Printf("[step %d]result: %+v\n", s+1, observation)

				functionCallOutputJSON, err := buildFunctionCallOutputJSON(toolCall.CallID, observation)
				if err != nil {
					return err
				}

				state.Transcript = append(state.Transcript, functionCallOutputJSON)
			}
		}
		if len(turn.ToolCalls) == 0 && turn.FinalAnswer != "" {
			state.History = append(state.History, step)
			return nil
		}
		state.History = append(state.History, step)
	}
	return fmt.Errorf("max steps reached")

}
