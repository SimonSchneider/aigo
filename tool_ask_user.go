package aigo

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

type ToolAskUser struct{}

func (t ToolAskUser) Name() string {
	return "User Prompt"
}

func (t ToolAskUser) Description() string {
	return "Interact with the user, ask a question or make a request, action input is the question or request. 'What is your name?' or 'How old are you?' or 'Can you help me complete this captcha?'"
}

func (t ToolAskUser) Exec(ctx context.Context, input string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("  LLM  > %s\n  User > ", input)
	text, _ := reader.ReadString('\n')
	return text, nil
}
