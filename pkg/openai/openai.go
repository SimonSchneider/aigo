package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/SimonSchneider/aigo/pkg/srv"
	"io"
	"net/http"
)

type Config struct {
	APIKey       string
	Organization string
}

type Client struct {
	client *http.Client
	config Config
}

func New(cfg Config, client *http.Client) *Client {
	mws := []srv.ClientMiddleware{
		srv.WithHeader("Authorization", "Bearer "+cfg.APIKey),
	}
	if cfg.Organization != "" {
		mws = append(mws, srv.WithHeader("OpenAI-Organization", cfg.Organization))
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &Client{
		client: srv.ClientWith(client, mws...),
		config: cfg,
	}
}

type Model struct {
	ID         string `json:"id"`
	OwnedBy    string `json:"owned_by"`
	Permission []any  `json:"permission"`
}

func (c *Client) Models(ctx context.Context) ([]Model, error) {
	resp, err := c.request(ctx, http.MethodGet, "https://api.openai.com/v1/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	defer resp.Body.Close()
	var body struct {
		Data []Model `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return body.Data, nil
}

const (
	RoleUser      = "user"
	RoleSystem    = "system"
	RoleAssistant = "assistant"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c ChatCompletionMessage) String() string {
	return fmt.Sprintf("%s: %s", c.Role, c.Content)
}

type ChatCompletionRequest struct {
	Model       string                  `json:"model"`
	Messages    []ChatCompletionMessage `json:"messages"`
	Temperature float64                 `json:"temperature,omitempty"`
	TopP        float64                 `json:"top_p,omitempty"`
	N           int                     `json:"n,omitempty"`
	Stop        []string                `json:"stop,omitempty"`
}

type ChatCompletionChoice struct {
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Choices []ChatCompletionChoice
}

func (c *Client) ChatCompletion(ctx context.Context, req ChatCompletionRequest, chatResp *ChatCompletionResponse) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	resp, err := c.request(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

func (c *Client) request(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}
	return resp, nil
}
