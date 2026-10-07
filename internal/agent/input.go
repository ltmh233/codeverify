package agent

import (
	"codeverify/internal/core"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

func buildInput(state *core.State) responses.ResponseInputParam {
	input := responses.ResponseInputParam{}
	for _, raw := range state.Transcript {
		input = append(input, param.Override[responses.ResponseInputItemUnionParam](raw))
	}
	return input
}

func buildFunctionCallOutputJSON(callID string, observation core.Observation) (json.RawMessage, error) {
	formatOutput := ""
	if observation.Success == true {
		formatOutput = fmt.Sprintf("[function call succeed, Output:] %s", observation.Output)
	} else {
		formatOutput = fmt.Sprintf("[function call falied, Output:] %s", observation.Output)
	}
	functionCallOutputJSON, err := responses.ResponseInputItemUnionParam{
		OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
			CallID: openai.String(callID),
			Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
				OfString: openai.String(formatOutput),
			},
		},
	}.MarshalJSON()
	if err != nil {
		return json.RawMessage{}, err
	}
	return functionCallOutputJSON, nil
}
