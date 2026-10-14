package tool

import (
	"codeverify/internal/analysis"
	"codeverify/internal/core"
	"context"
	"encoding/json"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type GetCallTool struct{}

func (t *GetCallTool) Name() core.ToolName {
	return core.ToolGetCall
}

func (t *GetCallTool) Execute(ctx context.Context, runCtx core.RunContext, rawARgs string) (core.Observation, error) {
	var args struct {
		Symbol string `json:"symbol"`
		Depth  int    `json:"depth"`
	}

	err := json.Unmarshal([]byte(rawARgs), &args)
	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, nil
	}

	if args.Symbol == "" {
		return core.Observation{
			Success: false,
			Output:  "have not pass argument symbol",
		}, nil
	}

	if args.Depth == 0 {
		return core.Observation{
			Success: false,
			Output:  "have not pass argument depth",
		}, nil
	}

	if args.Depth < 1 || args.Depth > 5 {
		return core.Observation{
			Success: false,
			Output:  "argument depth exceed range",
		}, nil
	}

	output, err := analysis.GetCodeContext(runCtx.Analysis.Functions, runCtx.Analysis.Graph, args.Symbol, args.Depth)
	if err != nil {
		return core.Observation{
			Success: false,
			Output:  err.Error(),
		}, nil
	}

	return core.Observation{
		Success: true,
		Output:  output,
	}, nil
}

func (t *GetCallTool) Spec() *responses.ToolUnionParam {
	return &responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"symbol": map[string]any{
						"type":        "string",
						"description": "Exact SymbolID of the target function or method, as listed in the repository overview.",
					},
					"depth": map[string]any{
						"type":        "integer",
						"description": "Maximum number of call edges to traverse from the target function. Start with 1 or 2 for focused analysis, and increase when deeper dependencies are relevant.",
						"minimum":     1,
						"maximum":     5,
					},
				},
				"required":             []string{"symbol", "depth"},
				"additionalProperties": false,
			},
			Name: "get_call",
			Description: openai.String(`Retrieve the code context for a specific function or method.
Returns the target function's signature, source location,
related functions within the specified call depth,
and their static call relationships.

Use this tool to understand function dependencies,
investigate potential bugs, and identify relevant code
before generating verification tests.

The symbol must be an exact SymbolID provided by the
repository overview. Do not invent or modify SymbolIDs.

This tool provides structural information only.
Read the source files to inspect actual implementations.`,
			)},
	}
}
