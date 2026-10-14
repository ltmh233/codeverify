package agent

import (
	"codeverify/internal/core"
	"codeverify/internal/tool"
	"context"
	"fmt"

	"github.com/openai/openai-go/v3/responses"
)

func Run(ctx context.Context, a Agent, runCtx core.RunContext, tools *tool.Registry) error {
	const maxSteps = 100

	for s := range maxSteps {
		turnStart := len(runCtx.State.Transcript)
		turn, err := a.Decide(ctx, runCtx.State, runCtx, *tools)
		if err != nil {
			return err
		}
		step := core.Step{
			Turn:        turn,
			ToolResults: []core.ToolResult{},
		}
		for _, toolCall := range step.Turn.ToolCalls {
			fmt.Printf("[step %d]toolCall: %+v\n", s+1, toolCall.Name)

			t, ok := tools.GetTool(string(toolCall.Name))
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

				runCtx.State.Transcript = append(runCtx.State.Transcript, functionCallOutputJSON)
			} else {

				observation, err := t.Execute(ctx, runCtx, toolCall.RawArgs)
				if err != nil {
					return err
				}
				observation.Output = core.SanitizePath(observation.Output, runCtx.RepoAbs)
				if t.Name() == core.ToolRunTest {
					observation.Output = fmt.Sprintf("evidence_id: %s,\n %s", toolCall.CallID, observation.Output)
				}
				step.ToolResults = append(step.ToolResults, core.ToolResult{
					CallID:      toolCall.CallID,
					Observation: observation,
				})
				fmt.Printf("[step %d]result: %+v\n", s+1, observation)

				functionCallOutputJSON, err := buildFunctionCallOutputJSON(toolCall.CallID, observation)
				if err != nil {
					return err
				}

				runCtx.State.Transcript = append(runCtx.State.Transcript, functionCallOutputJSON)
			}
		}
		runCtx.State.History = append(runCtx.State.History, step)
		runCtx.State.TurnStarts = append(runCtx.State.TurnStarts, turnStart)
		if len(step.Turn.ToolCalls) == 0 && step.Turn.FinalAnswer != "" {
			if runCtx.State.HasSuspectedRisk() {
				message, err := BuildInputMessageJSON(
					"Verification is incomplete: unresolved risks remain. "+
						"Submit a verdict for each pending risk before finishing.",
					responses.EasyInputMessageRoleUser)
				if err != nil {
					return err
				}
				runCtx.State.Transcript = append(runCtx.State.Transcript, message)
			} else {
				return nil
			}
		}
	}
	return fmt.Errorf("max steps reached")

}
