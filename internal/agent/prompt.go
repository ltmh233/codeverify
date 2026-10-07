package agent

import (
	"encoding/json"

	"github.com/openai/openai-go/v3/responses"
)

func BuildInputMessageJSON(content string, role responses.EasyInputMessageRole) (json.RawMessage, error) {
	inputMessageJSON, err :=
		responses.ResponseInputItemParamOfMessage(content, role).MarshalJSON()
	if err != nil {
		return json.RawMessage{}, err
	}
	return inputMessageJSON, err
}
