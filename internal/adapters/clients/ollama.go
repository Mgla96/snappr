package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sashabaranov/go-openai"
)

func NewOllamaClient(cfg openai.ClientConfig) *OllamaClient {
	return &OllamaClient{
		cfg: cfg,
	}
}

type OllamaClient struct {
	cfg openai.ClientConfig
}

/*
	curl http://localhost:11434/api/chat -d '{
	  "model": "llama3.2",
	  "messages": [
	    { "role": "user", "content": "why is the sky blue?" }
	  ]
	}'
*/
func (c *OllamaClient) CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (response openai.ChatCompletionResponse, err error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/chat", c.cfg.BaseURL), bytes.NewBuffer(reqBody))
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("failed to read response: %v", err)
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return response, nil
}
