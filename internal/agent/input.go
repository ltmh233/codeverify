package agent

import (
	"cmp"
	"codeverify/internal/analysis"
	"codeverify/internal/core"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

func buildInput(state *core.State, runCtx core.RunContext) (responses.ResponseInputParam, error) {
	const keepTurns = 6
	input := responses.ResponseInputParam{}

	appendInput := func(transcript []json.RawMessage) {
		for _, raw := range transcript {
			input = append(input, param.Override[responses.ResponseInputItemUnionParam](raw))
		}
	}

	if len(state.TurnStarts) <= keepTurns {
		appendInput(state.Transcript)
		return input, nil
	}

	baseEnd := state.TurnStarts[0]
	appendInput(state.Transcript[:baseEnd])
	start := state.TurnStarts[len(state.TurnStarts)-keepTurns]
	appendInput(state.Transcript[start:])

	overview := analysis.GetRepoOverview(runCtx.Analysis.Functions)
	output, err := BuildInputMessageJSON(overview, responses.EasyInputMessageRoleUser)
	if err != nil {
		return responses.ResponseInputParam{}, err
	}
	appendInput([]json.RawMessage{output})

	snapshot := buildRiskSnapshot(state)
	output, err = BuildInputMessageJSON(snapshot, responses.EasyInputMessageRoleUser)
	if err != nil {
		return responses.ResponseInputParam{}, err
	}
	appendInput([]json.RawMessage{output})

	return input, nil
}

func buildFunctionCallOutputJSON(callID string, observation core.Observation) (json.RawMessage, error) {
	formatOutput := ""
	if observation.Success == true {
		formatOutput = fmt.Sprintf("[function call succeed, Output:] %s", truncateOutput(observation.Output, 12000))
	} else {
		formatOutput = fmt.Sprintf("[function call falied, Output:] %s", truncateOutput(observation.Output, 12000))
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

func buildRiskSnapshot(state *core.State) string {
	snapshot := strings.Builder{}
	fmt.Fprintln(&snapshot, "current verification state:")

	ids := []string{}
	for id := range state.Risks {
		ids = append(ids, id)
	}
	slices.SortFunc(ids, func(a string, b string) int {
		na, _ := strconv.ParseInt(strings.TrimPrefix(a, "R"), 10, 64)
		nb, _ := strconv.ParseInt(strings.TrimPrefix(b, "R"), 10, 64)
		return cmp.Compare(na, nb)
	})
	for _, id := range ids {
		risk := state.Risks[id]
		fmt.Fprintf(&snapshot, "%s:\n\tHypothesis: %s\n\tContract: %s\n\tStatus: %s\n\tEvidenceIDs: [%s]\n\tVerdictReason: %s\n\n", risk.ID, risk.Hypothesis, risk.Contract, risk.Status, strings.Join(risk.EvidenceIDs, ","), risk.VerdictReason)
	}

	fmt.Fprintln(&snapshot, "Continue verification based on the current state.")

	return snapshot.String()
}

func truncateOutput(s string, limit int) string {
	if limit <= 0 {
		return "[output omitted]"
	}

	r := []rune(s)

	if len(r) <= limit {
		return s
	}

	head := 2 * limit / 3
	tail := limit - head

	return string(r[:head]) + "\n...[output truncate]...\n" + string(r[len(r)-tail:])

}
