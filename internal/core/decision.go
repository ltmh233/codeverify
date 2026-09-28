package core

type DecisionType string

const (
	DecisionToolCall DecisionType = "tool_call"
	DecisionFinish   DecisionType = "finish"
)

type ToolName string

const (
	ToolReadFile  ToolName = "read_file"
	ToolWriteTest ToolName = "write_test"
	ToolRunTest   ToolName = "run_test"
)

type ToolCall struct {
	Name ToolName
	Args map[string]string
}

type Decision struct {
	Type     DecisionType
	ToolCall *ToolCall
}
