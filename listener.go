package aigo

import (
	"fmt"
	"github.com/SimonSchneider/aigo/pkg/cfmt"
	"github.com/SimonSchneider/aigo/prompts"
)

type Listener interface {
	OnQuestion(question string)
	OnResponse(response prompts.AgentResponse)
	OnObservation(observation string)
	OnPrompt(prompt string)
}

type ListenerKey string

const (
	ListenerKeyThought     ListenerKey = "Thought"
	ListenerKeyAction      ListenerKey = "Action"
	ListenerKeyActionInput ListenerKey = "ActionInput"
	ListenerKeyObservation ListenerKey = "Observation"
	ListenerKeyQuestion    ListenerKey = "Question"
)

var colorMapping = map[ListenerKey]cfmt.Color{
	ListenerKeyQuestion:    cfmt.Purple,
	ListenerKeyThought:     cfmt.Green,
	ListenerKeyAction:      cfmt.Red,
	ListenerKeyActionInput: cfmt.Orange,
	ListenerKeyObservation: cfmt.Blue,
}

type TerminalListener struct {
	Debug bool
}

func (t TerminalListener) OnQuestion(question string) {
	t.Print(ListenerKeyQuestion, question)
}

func (t TerminalListener) OnResponse(response prompts.AgentResponse) {
	t.Print(ListenerKeyThought, response.Thought)
	t.Print(ListenerKeyAction, response.Action)
	t.Print(ListenerKeyActionInput, response.ActionInput)
}

func (t TerminalListener) OnObservation(observation string) {
	t.Print(ListenerKeyObservation, observation)
}

func (t TerminalListener) Print(action ListenerKey, s string) {
	_, _ = cfmt.Print(colorMapping[action], string(action), s)
}

func (t TerminalListener) OnPrompt(prompt string) {
	if t.Debug {
		fmt.Printf("---- Debug Prompt: \n%s\n----\n", prompt)
	}
}

func NewTerminalListener() TerminalListener {
	return TerminalListener{}
}

type NoopListener struct{}

func (n NoopListener) OnQuestion(question string)                {}
func (n NoopListener) OnResponse(response prompts.AgentResponse) {}
func (n NoopListener) OnObservation(observation string)          {}
func (n NoopListener) OnPrompt(observation string)               {}
