package tools

import (
	"fmt"
	"strings"
)

type Params struct {
	ImagePath string
	Prompt    string
	Options   map[string]any
}

type Result struct {
	ImageBytes []byte
	ImageURL   string
	Metadata   map[string]string
}

type Tool interface {
	Name() string
	Description() string
	Execute(params Params) (*Result, error)
}

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) Register(t Tool) {
	r.tools[t.Name()] = t
}

func (r *Registry) Get(name string) (Tool, error) {
	t, ok := r.tools[name]
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	return t, nil
}

func (r *Registry) List() []Tool {
	var list []Tool
	for _, t := range r.tools {
		list = append(list, t)
	}
	return list
}

func (r *Registry) SystemPromptSection() string {
	var sb strings.Builder
	sb.WriteString("You can use the following tools:\n\n")
	for i, t := range r.List() {
		fmt.Fprintf(&sb, "%d. %s - %s\n", i+1, t.Name(), t.Description())
	}
	return sb.String()
}

func (r *Registry) EnabledToolNames() []string {
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}
