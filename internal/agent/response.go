package agent

import (
	"codeverify/internal/core"
	"encoding/json"
	"log"

	"github.com/openai/openai-go/v3/responses"
)

func parseResponse(res *responses.Response, state *core.State) (core.Turn, error) {
	turn := core.Turn{
		Reasonings:  []string{},
		ToolCalls:   []*core.ToolCall{},
		Messages:    []string{},
		FinalAnswer: "",
	}
	for _, output := range res.Output {
		switch output.Type {
		case "function_call":
			call := output.AsFunctionCall()
			state.Transcript = append(state.Transcript, json.RawMessage(call.RawJSON()))
			turn.ToolCalls = append(turn.ToolCalls, &core.ToolCall{
				Name:    core.ToolName(call.Name),
				RawArgs: call.Arguments,
				CallID:  call.CallID,
			})
		case "message":
			message := output.AsMessage()
			state.Transcript = append(state.Transcript, json.RawMessage(message.RawJSON()))
			if message.Phase == "final_answer" {
				finalAnswer := ""
				for _, content := range message.Content {
					finalAnswer = finalAnswer + content.Text
				}
				turn.FinalAnswer = finalAnswer
			}
		case "reasoning":
			reasoning := output.AsReasoning()
			state.Transcript = append(state.Transcript, json.RawMessage(reasoning.RawJSON()))
			for _, content := range reasoning.Content {
				if content.Type == "reasoning_text" {
					turn.Reasonings = append(turn.Reasonings, content.Text)
				}
			}
		default:
			log.Printf("unknown type of output.Type")
			continue
		}
	}

	return turn, nil
}
