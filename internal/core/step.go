package core

type Observation struct {
	Success bool
	Output  string
}

type ToolResult struct {
	CallID      string
	Observation Observation
}

type Turn struct {
	Reasonings  []string
	ToolCalls   []*ToolCall
	Messages    []string
	FinalAnswer string
}

type Step struct {
	Turn        Turn
	ToolResults []ToolResult
}
