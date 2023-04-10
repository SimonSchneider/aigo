package prompts

import (
	"context"
	"fmt"
)

type AgentPrompt struct {
	Question          string
	Tools             []Tool
	PreviousResponses []ChainLink
}

type Tool interface {
	Name() string
	Description() string
	Exec(ctx context.Context, input string) (string, error)
}

type ChainLink struct {
	AgentResponse
	ActionOutput string
}

type AgentResponse struct {
	Thought     string
	Action      string
	ActionInput string
}

func (a AgentResponse) String() string {
	return fmt.Sprintf("Thought: %s\nAction: %s\nActionInput: %s\n", a.Thought, a.Action, a.ActionInput)
}
