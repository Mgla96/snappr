package clients

import (
	"context"

	"github.com/Mgla96/snappr/internal/errors"
	"github.com/sashabaranov/go-openai"
)

type ModelType string

const (
	// GPT3_5Turbo0125 is the GPT-3.5-turbo-0125 model.
	GPT3_5Turbo0125            ModelType = "gpt-3.5-turbo-0125"
	GPT4_turbo                 ModelType = "gpt-4-turbo"
	ErrNoChatCompletionChoices           = errors.New("no chat completion choices returned")
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// openAIBackingClient is an interface for openai.Client
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . openAIBackingClient
type openAIBackingClient interface {
	CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (response openai.ChatCompletionResponse, err error)
}

type OpenAIClient struct {
	aiClient openAIBackingClient
}

// ModelToContextWindow returns the context window size for the given model.
func ModelToContextWindow(model ModelType) int {
	switch model {
	case GPT3_5Turbo0125:
		return 16385
	case GPT4_turbo:
		return 128000
	default:
		return -1
	}
}

func NewCustomOpenAIClient(authToken, baseURL string) *OpenAIClient {
	defaultConfig := openai.DefaultConfig(authToken)
	defaultConfig.BaseURL = baseURL
	client := openai.NewClientWithConfig(defaultConfig)
	return &OpenAIClient{
		aiClient: client,
	}
}

// NewOpenAIClient creates a new instance of the OpenAIClient.
func NewOpenAIClient(apiKey string) *OpenAIClient {
	defaultConfig := openai.DefaultConfig(apiKey)
	// defaultConfig.BaseURL = "http://localhost:11434"
	client := openai.NewClientWithConfig(defaultConfig)
	return &OpenAIClient{
		aiClient: client,
	}
}

// GenerateChatCompletion generates a chat completion based on the provided messages.
//
// Parameters:
//   - ctx: The context for the API request.
//   - messages: The list of messages to use as input for the chat completion.
//
// Returns:
//   - The generated chat completion.
//   - An error if any occurred during the API request.
func (oc *OpenAIClient) GenerateChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage, model ModelType) (string, error) {
	req := openai.ChatCompletionRequest{
		Model:    string(model),
		Messages: messages,
	}

	resp, err := oc.aiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", ErrNoChatCompletionChoices
	}

	return resp.Choices[0].Message.Content, nil
}
