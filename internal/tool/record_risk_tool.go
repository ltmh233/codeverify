package tool

import (
	"codeverify/internal/core"
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type RecordRiskTool struct{}

func (t *RecordRiskTool) Name() core.ToolName {
	return core.ToolRecordRisk
}

func (t *RecordRiskTool) Execute(ctx context.Context, runCtx core.RunContext, rawArgs string) (core.Observation, error) {
	var args struct {
		Hypothesis string `json:"hypothesis"`
		Contract   string `json:"contract"`
	}
	err := json.Unmarshal([]byte(rawArgs), &args)
	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, nil
	}

	if args.Hypothesis == "" {
		return core.Observation{
			Success: false,
			Output:  "have not pass argument hypothesis",
		}, nil
	}

	if args.Contract == "" {
		return core.Observation{
			Success: false,
			Output:  "have not pass argument contract",
		}, nil
	}
	risk := runCtx.State.AddRisk(args.Hypothesis, args.Contract)

	return core.Observation{
		Success: true,
		Output:  fmt.Sprintf("riskID:%s ,riskStatus:%s", risk.ID, risk.Status),
	}, nil
}

func (t *RecordRiskTool) Spec() *responses.ToolUnionParam {
	return &responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"hypothesis": map[string]any{
						"type":        "string",
						"description": "A specific, testable hypothesis describing the potential defect.",
					},
					"contract": map[string]any{
						"type":        "string",
						"description": "The expected behavior and its supporting source in the repository.",
					},
				},
				"required":             []string{"hypothesis", "contract"},
				"additionalProperties": false,
			},
			Name: string(core.ToolRecordRisk),
			Description: openai.String(
				"Record a testable potential defect with its expected behavior " +
					"and supporting repository evidence before verification. " +
					"Returns a unique risk ID with suspected status.",
			),
		},
	}
}
