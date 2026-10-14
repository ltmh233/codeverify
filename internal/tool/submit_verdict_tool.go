package tool

import (
	"codeverify/internal/core"
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type SubmitVerdictTool struct {
}

func (t *SubmitVerdictTool) Name() core.ToolName {
	return core.ToolSubmitVerdict
}

func (t *SubmitVerdictTool) Execute(ctx context.Context, runCtx core.RunContext, rawArgs string) (core.Observation, error) {
	var args struct {
		RiskID      string   `json:"risk_id"`
		Status      string   `json:"status"`
		EvidenceIDs []string `json:"evidence_ids"`
		Reason      string   `json:"reason"`
	}
	err := json.Unmarshal([]byte(rawArgs), &args)
	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, nil
	}

	err = runCtx.State.SubmitVerdict(args.RiskID, args.Status, args.EvidenceIDs, args.Reason)
	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, nil
	}

	return core.Observation{
		Success: true,
		Output:  fmt.Sprintf("update risk %s status to %s", args.RiskID, args.Status),
	}, nil
}

func (t *SubmitVerdictTool) Spec() *responses.ToolUnionParam {
	return &responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"risk_id": map[string]any{
						"type":        "string",
						"description": "ID of the risk being evaluated.",
					},
					"status": map[string]any{
						"type": "string",
						"enum": []string{
							string(core.RiskConfirmed),
							string(core.RiskRejected),
							string(core.RiskInconclusive),
						},
						"description": "Status of the risk by verdict want to set",
					},
					"evidence_ids": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
						"description": "Evidence_IDs the verdict supported by",
					},
					"reason": map[string]any{
						"type":        "string",
						"description": "Explanation of the verdict.",
					},
					"additionalProperties": false,
				},
			},
			Name: string(core.ToolSubmitVerdict),
			Description: openai.String(
				"Submit a verification verdict for an existing risk. " +
					"Provide the risk ID, verdict status, supporting execution " +
					"evidence IDs, and an explanation of how the evidence " +
					"supports the conclusion. A failing test alone does not " +
					"automatically confirm a defect.",
			),
		},
	}
}
