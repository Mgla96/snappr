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
	ollamaReq, err := convertToOllamaRequest(request)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("failed to convert request: %v", err)
	}
	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("failed to marshal request: %v", err)
	}

	fmt.Printf("request: %s\n", reqBody)

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

	var ollamaResponse ollamaResponse
	fmt.Printf("response: %s\n", respBody)

	err = json.Unmarshal(respBody, &ollamaResponse)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// convert ollamaResponse to openai.ChatCompletionResponse
	return convertOllamaRespToOpenAIResponse(ollamaResponse), nil
}

type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}
type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Model              string        `json:"model"`
	CreatedAt          string        `json:"created_at"`
	Message            ollamaMessage `json:"message"`
	Done               bool          `json:"done"`
	TotalDuration      float64       `json:"total_duration"`
	PromptEvalCount    int           `json:"prompt_eval_count"`
	PromptEvalDuration float64       `json:"prompt_eval_duration"`
	EvalCount          int           `json:"eval_count"`
	EvalDuration       float64       `json:"eval_duration"`
}

func convertToOllamaRequest(request openai.ChatCompletionRequest) (ollamaRequest, error) {
	ollamaRequest := ollamaRequest{
		Model:    request.Model,
		Messages: []ollamaMessage{},
		Stream:   request.Stream,
	}

	for _, message := range request.Messages {
		ollamaMessage := ollamaMessage{
			Role:    message.Role,
			Content: message.Content,
		}

		ollamaRequest.Messages = append(ollamaRequest.Messages, ollamaMessage)
	}

	return ollamaRequest, nil
}

func convertOllamaRespToOpenAIResponse(response ollamaResponse) openai.ChatCompletionResponse {
	return openai.ChatCompletionResponse{
		Model: response.Model,
		Choices: []openai.ChatCompletionChoice{
			{
				Message: openai.ChatCompletionMessage{
					Role:    response.Message.Role,
					Content: response.Message.Content,
				},
			},
		},
	}
}
