package analysis

type CallEdge struct {
	Caller   string
	Callee   string
	CallFile string
	CallLine int
}

type CallGraph struct {
	Edges []CallEdge
}

func (g *CallGraph) Reachable(startSymbol string, maxDepth int) ([]CallEdge, []string) {
	relatedFunctions := []string{}
	adj := map[string][]CallEdge{}
	for _, edge := range g.Edges {
		adj[edge.Caller] = append(adj[edge.Caller], edge)
	}
	visit := map[string]bool{}
	type node struct {
		Caller string
		depth  int
	}

	queue := []node{}
	queue = append(queue, node{
		Caller: startSymbol,
		depth:  0,
	})

	edges := []CallEdge{}

	for len(queue) != 0 && queue[0].depth < maxDepth {
		front := queue[0]
		queue = queue[1:]
		if visit[front.Caller] == true {
			continue
		}
		relatedFunctions = append(relatedFunctions, front.Caller)
		visit[front.Caller] = true
		for _, edge := range adj[front.Caller] {
			queue = append(queue, node{
				Caller: edge.Callee,
				depth:  front.depth + 1,
			})
			edges = append(edges, edge)
		}
	}

	return edges, relatedFunctions
}
