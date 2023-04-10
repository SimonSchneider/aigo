package aigo

import (
	"embed"
	"fmt"
	"github.com/SimonSchneider/aigo/prompts"
	"io"
	"text/template"
)

//go:embed prompts/*
var promptFS embed.FS

const (
	PromptAgent = "agent.tmpl"
	Prompt2     = "prompt2.tmpl"
)

func ParsePromptTemplates(t *template.Template) (*template.Template, error) {
	return t.ParseFS(promptFS, "prompts/*.tmpl")
}

type PromptGenerator struct {
	templates *template.Template
}

func NewPromptGenerator() (*PromptGenerator, error) {
	t, err := ParsePromptTemplates(template.New("prompt"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse prompt templates: %w", err)
	}
	return &PromptGenerator{templates: t}, nil
}

func MustNewPromptGenerator() *PromptGenerator {
	pg, err := NewPromptGenerator()
	if err != nil {
		panic(err)
	}
	return pg
}

func (p *PromptGenerator) GenerateAgentPrompt(w io.Writer, prompt *prompts.AgentPrompt) error {
	return p.executeTemplate(w, PromptAgent, prompt)
}

func (p *PromptGenerator) executeTemplate(w io.Writer, name string, data any) error {
	if err := p.templates.ExecuteTemplate(w, name, data); err != nil {
		return fmt.Errorf("failed to execute template '%s': %w", name, err)
	}
	return nil
}
