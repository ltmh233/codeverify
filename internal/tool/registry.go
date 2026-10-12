package tool

import "github.com/openai/openai-go/v3/responses"

type Registry struct {
	tools map[string]Tool
}

func (r *Registry) Specs() []responses.ToolUnionParam {
	tools := []responses.ToolUnionParam{}
	for _, tool := range r.tools {
		tools = append(tools, *tool.Spec())
	}
	return tools
}

func (r *Registry) Register(tool Tool) {
	if r.tools == nil {
		r.tools = make(map[string]Tool)
	}
	r.tools[string(tool.Name())] = tool
}

func (r *Registry) GetTool(name string) (Tool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}
