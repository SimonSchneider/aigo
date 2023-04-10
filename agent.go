package aigo

import (
	"context"
	"fmt"
	"github.com/SimonSchneider/aigo/pkg/openai"
	"github.com/SimonSchneider/aigo/prompts"
	"regexp"
	"strings"
)

const (
	ActionNone        = "None"
	ActionFinalAnswer = "FinalAnswer"
)

type OpenAIClient interface {
	ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest, o *openai.ChatCompletionResponse) error
}

type Agent struct {
	pg       *PromptGenerator
	oai      OpenAIClient
	listener Listener
	model    string
}

func NewAgent(pg *PromptGenerator, oai OpenAIClient, model string, listener Listener) *Agent {
	if listener == nil {
		listener = NoopListener{}
	}
	return &Agent{pg: pg, oai: oai, model: model, listener: listener}
}

type AgentQuestion struct {
	Question    string
	MaxRequests int
	Tools       []prompts.Tool
}

type AgentResponse struct {
	Answer string
}

func (a *Agent) Ask(ctx context.Context, question AgentQuestion) (AgentResponse, error) {
	toolExecutor := NewToolExecutor(question.Tools)
	prompt := prompts.AgentPrompt{
		Question:          question.Question,
		Tools:             question.Tools,
		PreviousResponses: make([]prompts.ChainLink, 0, question.MaxRequests),
	}
	a.listener.OnQuestion(question.Question)
	for i := 0; i < question.MaxRequests; i++ {
		var res prompts.AgentResponse
		if err := a.Request(ctx, &prompt, &res); err != nil {
			return AgentResponse{}, fmt.Errorf("failed to request agent: %w", err)
		}
		a.listener.OnResponse(res)
		prompt.PreviousResponses = append(prompt.PreviousResponses, prompts.ChainLink{AgentResponse: res})
		if res.Action == ActionNone {
			return AgentResponse{
				Answer: res.Thought,
			}, nil
		} else if res.Action == ActionFinalAnswer {
			return AgentResponse{
				Answer: res.ActionInput,
			}, nil
		}
		toolRes := toolExecutor.Execute(ctx, res.Action, res.ActionInput)
		a.listener.OnObservation(toolRes)
		prompt.PreviousResponses[len(prompt.PreviousResponses)-1].ActionOutput = toolRes
	}
	return AgentResponse{}, nil
}

func (a *Agent) Request(ctx context.Context, prompt *prompts.AgentPrompt, r *prompts.AgentResponse) error {
	var buf strings.Builder
	if err := a.pg.GenerateAgentPrompt(&buf, prompt); err != nil {
		return fmt.Errorf("failed to generate prompt: %w", err)
	}
	promptStr := buf.String()
	a.listener.OnPrompt(promptStr)
	chatReq := openai.ChatCompletionRequest{
		Model: a.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.RoleUser, Content: promptStr},
		},
		Stop: []string{"Observation:"},
	}
	var resp openai.ChatCompletionResponse
	if err := a.oai.ChatCompletion(ctx, chatReq, &resp); err != nil {
		return fmt.Errorf("failed to complete prompt: %w", err)
	}
	if (len(resp.Choices)) == 0 {
		return fmt.Errorf("no response from agent")
	}
	respMsg := resp.Choices[0].Message.Content
	if err := ParseAgentResponse(r, respMsg); err != nil {
		return fmt.Errorf("failed to parse agent response: %w", err)
	}
	return nil
}

var actionPattern = regexp.MustCompile(`(?si)Thought:(.*?)Action:(.*?)(?:Action\s?Input:(.*))?$`)
var finalAnswerPat = regexp.MustCompile(`(?si)Thought:(.*?)Final\s?Answer:(.*)$`)

const whiteSpace = "\n \t"

func ParseAgentResponse(r *prompts.AgentResponse, respMsg string) error {
	matches := actionPattern.FindStringSubmatch(respMsg)
	for i := range matches {
		matches[i] = strings.Trim(matches[i], whiteSpace)
	}
	nMatches := len(matches)
	if nMatches == 0 {
		matches = finalAnswerPat.FindStringSubmatch(respMsg)
		for i := range matches {
			matches[i] = strings.Trim(matches[i], whiteSpace)
		}
		if len(matches) == 3 {
			r.Thought = matches[1]
			r.Action = ActionFinalAnswer
			r.ActionInput = matches[2]
			return nil
		}
		return fmt.Errorf("invalid response from agent: %s", respMsg)
	}
	if nMatches < 3 || nMatches > 4 {
		return fmt.Errorf("invalid response from agent: %s", respMsg)
	}
	r.Thought = matches[1]
	if nMatches == 4 {
		r.Action = matches[2]
		r.ActionInput = matches[3]
	} else if nMatches == 3 && PrefixFold(r.Action, "none") {
		r.Action = ActionNone
	} else {
		return fmt.Errorf("invalid response from agent: %s", respMsg)
	}
	return nil
}

func PrefixFold(s, prefix string) bool {
	return strings.EqualFold(s[:len(prefix)], prefix)
}

type ToolExecutor map[string]prompts.Tool

func NewToolExecutor(tools []prompts.Tool) ToolExecutor {
	m := make(map[string]prompts.Tool, len(tools))
	for _, tool := range tools {
		m[tool.Name()] = tool
	}
	return m
}

func (t *ToolExecutor) Execute(ctx context.Context, name string, input string) string {
	tool := (*t)[name]
	if tool == nil {
		return "unknown tool"
	}
	res, err := tool.Exec(ctx, input)
	if err != nil {
		return fmt.Sprintf("error: %s", err)
	}
	return strings.Trim(res, whiteSpace)
}
