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

// NewOllamaClient constructs a new OllamaClient with the provided configuration.
func NewOllamaClient(cfg openai.ClientConfig) *OllamaClient {
	return &OllamaClient{
		cfg: cfg,
	}
}

// OllamaClient manages chat completions for the Ollama model.
type OllamaClient struct {
	cfg openai.ClientConfig
}

// CreateChatCompletion sends a request to the Ollama server and returns the response.
func (c *OllamaClient) CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (response openai.ChatCompletionResponse, err error) {...}

// Converts an OpenAI request format to Ollama's required format.
func convertToOllamaRequest(request openai.ChatCompletionRequest) ollamaRequest {...}

// Converts Ollama's response back into OpenAI's response format.
func convertOllamaRespToOpenAIResponse(response ollamaResponse) openai.ChatCompletionResponse {...}
