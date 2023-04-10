package aigo

import (
	"github.com/SimonSchneider/aigo/prompts"
	"strings"
	"testing"
)

func TestReadFromRegex(t *testing.T) {
	// table test for the Parse agent Response
	tests := []struct {
		name    string
		str     string
		want    prompts.AgentResponse
		wantErr bool
	}{
		{
			name: "valid",
			str:  "Thought: test\nAction: test\nhello\nActionInput: test\n",
			want: prompts.AgentResponse{
				Thought:     "test",
				Action:      "test\nhello",
				ActionInput: "test",
			},
			wantErr: false,
		},
		{
			name: "lower case match",
			str:  "Thought: test\naction: test\nActionInput: test",
			want: prompts.AgentResponse{
				Thought:     "test",
				Action:      "test",
				ActionInput: "test",
			},
			wantErr: false,
		},
		{
			name: "real world",
			str: `Thought: I want to use one of the available tools to gather some new information.
Action: pizza-delivery
Action Input: [mushroom; 456 oak avenue]
`,
			want: prompts.AgentResponse{
				Thought:     "I want to use one of the available tools to gather some new information.",
				Action:      "pizza-delivery",
				ActionInput: `[mushroom; 456 oak avenue]`,
			},
		},
		{
			name: "without action input",
			str:  "Thought: I am an AI language model and do not eat, so I will not be ordering a pizza.\nAction: None",
			want: prompts.AgentResponse{
				Thought:     "I am an AI language model and do not eat, so I will not be ordering a pizza.",
				Action:      "None",
				ActionInput: "",
			},
		},
		{
			name: "final answer",
			str:  "Thought: Since I am an AI language model, I do not have personal preferences or desires. \nFinal Answer: As an AI language model, I do not want to do anything.",
			want: prompts.AgentResponse{
				Thought:     "Since I am an AI language model, I do not have personal preferences or desires.",
				Action:      "FinalAnswer",
				ActionInput: "As an AI language model, I do not want to do anything.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got prompts.AgentResponse
			err := ParseAgentResponse(&got, tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAgentResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.EqualFold(got.Thought, tt.want.Thought) {
				t.Errorf("parseAgentResponse(Thought) \ngot  = %+v\nwant = %+v", got.Thought, tt.want.Thought)
			}
			if !strings.EqualFold(got.Action, tt.want.Action) {
				t.Errorf("parseAgentResponse() got = %v, want %v", got, tt.want)
			}
			if !strings.EqualFold(got.ActionInput, tt.want.ActionInput) {
				t.Errorf("parseAgentResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
