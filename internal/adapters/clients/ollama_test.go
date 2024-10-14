package clients

import (
	"context"
	"reflect"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestOllamaClient_CreateChatCompletion(t *testing.T) {
	type fields struct {
		cfg openai.ClientConfig
	}
	type args struct {
		ctx     context.Context
		request openai.ChatCompletionRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse openai.ChatCompletionResponse
		wantErr      bool
	}{
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &OllamaClient{
				cfg: tt.fields.cfg,
			}
			gotResponse, err := c.CreateChatCompletion(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("OllamaClient.CreateChatCompletion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("OllamaClient.CreateChatCompletion() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestNewOllamaClient(t *testing.T) {
	type args struct {
		cfg openai.ClientConfig
	}
	tests := []struct {
		name string
		args args
		want *OllamaClient
	}{
		{
			name: "create new ollama client",
			args: args{
				cfg: openai.ClientConfig{
					BaseURL: "http://localhost:11434",
				},
			},
			want: &OllamaClient{
				cfg: openai.ClientConfig{
					BaseURL: "http://localhost:11434",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOllamaClient(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOllamaClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToOllamaRequest(t *testing.T) {
	type args struct {
		request openai.ChatCompletionRequest
	}
	tests := []struct {
		name string
		args args
		want ollamaRequest
	}{
		{
			name: "successful",
			args: args{
				request: openai.ChatCompletionRequest{
					Model: "llama3.2",
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    "user",
							Content: "why is the sky blue?",
						},
					},
					Stream: false,
				},
			},
			want: ollamaRequest{
				Model: "llama3.2",
				Messages: []ollamaMessage{
					{
						Role:    "user",
						Content: "why is the sky blue?",
					},
				},
				Stream: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertToOllamaRequest(tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToOllamaRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertOllamaRespToOpenAIResponse(t *testing.T) {
	type args struct {
		response ollamaResponse
	}
	tests := []struct {
		name string
		args args
		want openai.ChatCompletionResponse
	}{
		{
			name: "successful",
			args: args{
				response: ollamaResponse{
					Model:     "llama3.2",
					CreatedAt: "2022-01-01T00:00:00Z",
					Message: ollamaMessage{
						Role:    "user",
						Content: "foobar",
					},
				},
			},
			want: openai.ChatCompletionResponse{
				Model: "llama3.2",
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Role:    "user",
							Content: "foobar",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertOllamaRespToOpenAIResponse(tt.args.response); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertOllamaRespToOpenAIResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
