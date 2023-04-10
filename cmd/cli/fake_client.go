package main

import (
	"context"
	"fmt"
	"github.com/SimonSchneider/aigo/pkg/openai"
	"strings"
)

type FakeClient struct{}

func (f FakeClient) Models(ctx context.Context) ([]openai.Model, error) {
	models := []openai.Model{
		{
			ID:      "gpt-3.5-turbo",
			OwnedBy: "fakers",
		},
	}
	return models, nil
}

func (f FakeClient) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest, o *openai.ChatCompletionResponse) error {
	const (
		qImpossible = "impossible"
		qEasy       = "easy"
		qHard       = "hard"
	)
	lastMsg := strings.Split(strings.Trim(req.Messages[len(req.Messages)-1].Content, " \n\t"), "\n")

	lastPrompt := lastMsg[len(lastMsg)-1]
	if question, ok := strings.CutPrefix(lastPrompt, "Question: "); ok {
		if question == qImpossible {
			o.Choices = f.answerWith("Thought: can't answer this Action: None")
		} else if question == qEasy {
			o.Choices = f.answerWith("Thought: easy I know this Final Answer: because I say so")
		} else if question == qHard {
			o.Choices = f.answerWith("Thought: hard I don't know this Action: User Prompt Action Input: whats the answer")
		}
	} else if observation, ok := strings.CutPrefix(lastPrompt, "Observation: "); ok {
		o.Choices = f.answerWith(fmt.Sprintf("Thought: now I know FinalAnswer: %s", observation))
	} else {
		o.Choices = f.answerWith(fmt.Sprintf("Thought: I don't know, I only know [%s] Action: None", strings.Join([]string{qImpossible, qEasy, qHard}, ", ")))
	}
	return nil
}

func (f FakeClient) answerWith(content string) []openai.ChatCompletionChoice {
	return []openai.ChatCompletionChoice{
		{Index: 0, Message: openai.ChatCompletionMessage{Role: "assistant", Content: content}, FinishReason: "done"},
	}
}
