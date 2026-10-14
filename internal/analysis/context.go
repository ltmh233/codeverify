package analysis

import (
	"fmt"
	"strings"
)

type Analysis struct {
	Functions []FunctionInfo
	Graph     *CallGraph
}

func displayName(info FunctionInfo) string {
	if info.Receiver != "" {
		return fmt.Sprintf("(%s).%s", info.Receiver, info.Name)
	}
	return info.Name
}

func formatFunction(info FunctionInfo) string {
	sig := strings.TrimPrefix(info.Signature, "func")

	if info.Receiver != "" {
		return fmt.Sprintf("func (%s) %s%s",
			info.Receiver, info.Name, sig)
	}

	return fmt.Sprintf("func %s%s", info.Name, sig)
}

func GetCodeContext(functions []FunctionInfo, graph *CallGraph, startSymbol string, depth int) (string, error) {
	index := make(map[string]FunctionInfo)
	for _, info := range functions {
		index[info.Symbol] = info
	}

	edges, relatedFunctions := graph.Reachable(startSymbol, 100)

	builder := strings.Builder{}

	info, ok := index[startSymbol]
	if !ok {
		return "", fmt.Errorf("%s function not exists", startSymbol)
	}
	fmt.Fprintf(&builder, "Target:\n\t%s\n\n", formatFunction(info))

	fmt.Fprintf(&builder, "Related Functions:\n")
	for _, function := range relatedFunctions {
		if function == startSymbol {
			continue
		}
		info, ok := index[function]
		if !ok {
			return "", fmt.Errorf("%s function not exists", function)
		}
		fmt.Fprintf(&builder, "\t%s\n\tFile: %s:%d\n\n", formatFunction(info), info.File, info.Line)
	}

	fmt.Fprintf(&builder, "Call Relations:\n")
	for _, edge := range edges {
		info1, ok := index[edge.Caller]
		if !ok {
			return "", fmt.Errorf("%s function not exists", edge.Caller)
		}
		info2, ok := index[edge.Callee]
		if !ok {
			return "", fmt.Errorf("%s function not exists", edge.Callee)
		}
		fmt.Fprintf(&builder, "\t%s -> %s\n", displayName(info1), displayName(info2))
	}
	return builder.String(), nil
}

func GetRepoOverview(functions []FunctionInfo) string {
	pkgs := map[string][]FunctionInfo{}
	for _, function := range functions {
		pkgs[function.Package] = append(pkgs[function.Package], function)
	}

	builder := strings.Builder{}

	fmt.Fprintf(&builder, "Repository Overview:\n")
	for pkg, infos := range pkgs {
		fmt.Fprintf(&builder, "\nPackage: %s\n\n", pkg)
		for _, info := range infos {
			fmt.Fprintf(&builder, "\tFunction: %s\n", formatFunction(info))
			fmt.Fprintf(&builder, "\t\tSymbol: %s\n", info.Symbol)
			fmt.Fprintf(&builder, "\t\tLocation: %s:%d\n\n", info.File, info.Line)
		}
	}

	return builder.String()
}
