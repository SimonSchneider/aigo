package aigo

import (
	"context"
	"fmt"
	"github.com/SimonSchneider/aigo/pkg/openai"
	"sort"
	"strings"
)

type Models []openai.Model

type OpenAIModelClient interface {
	Models(ctx context.Context) ([]openai.Model, error)
}

func GetModels(ctx context.Context, client OpenAIModelClient) (Models, error) {
	models, err := client.Models(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}
	Models(models).sort()
	return models, nil
}

func (m Models) sort() {
	sort.Slice(m, func(i, j int) bool {
		return m[i].ID < m[j].ID
	})
}

func (m Models) Select(modelPriority ...string) (string, error) {
	for _, model := range modelPriority {
		_, ok := sort.Find(len(m), func(i int) int {
			return strings.Compare(model, m[i].ID)
		})
		if ok {
			return model, nil
		}
	}
	availableModels := make([]string, len(m))
	for i, model := range m {
		availableModels[i] = model.ID
	}
	return "", fmt.Errorf("none of the given models are available, available models: %v", availableModels)
}
