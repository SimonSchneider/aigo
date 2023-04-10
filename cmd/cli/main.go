package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/SimonSchneider/aigo"
	"github.com/SimonSchneider/aigo/pkg/openai"
	"github.com/SimonSchneider/aigo/prompts"
	"os"
)

func loadCfg() (cfg Cfg, err error) {
	flag.BoolVar(&cfg.Debug, "debug", false, "enable debug mode")
	flag.BoolVar(&cfg.Fake, "fake", false, "use fake openai client")
	flag.StringVar(&cfg.ApiKey, "openai-key", "", "openai api key")
	flag.StringVar(&cfg.Question, "question", "", "question to ask")
	flag.Parse()
	return
}

func getModel(ctx context.Context, client aigo.OpenAIModelClient, prio ...string) (string, error) {
	models, err := aigo.GetModels(ctx, client)
	if err != nil {
		return "", fmt.Errorf("failed to get models: %w", err)
	}
	model, err := models.Select(prio...)
	if err != nil {
		return "", fmt.Errorf("failed to select high perf model: %w", err)
	}
	return model, nil
}

type OpenAIClient interface {
	aigo.OpenAIClient
	aigo.OpenAIModelClient
}

func run(ctx context.Context, cfg Cfg) error {
	var oai OpenAIClient
	if cfg.Fake {
		oai = FakeClient{}
	} else {
		oai = openai.New(openai.Config{
			APIKey: cfg.ApiKey,
		}, nil)
	}
	model, err := getModel(ctx, oai, "gpt-4-32k", "gpt-4", "gpt-3.5-turbo")
	if err != nil {
		return err
	}
	pg := aigo.MustNewPromptGenerator()
	listener := aigo.TerminalListener{Debug: cfg.Debug}
	agent := aigo.NewAgent(pg, oai, model, listener)
	_, err = agent.Ask(ctx, aigo.AgentQuestion{
		Question:    cfg.Question,
		MaxRequests: 10,
		Tools: []prompts.Tool{
			aigo.ToolAskUser{},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to ask question: %w", err)
	}
	return nil
}

type Cfg struct {
	ApiKey   string
	Debug    bool
	Fake     bool
	Question string
}

func main() {
	cfg, err := loadCfg()
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		os.Exit(1)
	}
	if err := run(context.TODO(), cfg); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
