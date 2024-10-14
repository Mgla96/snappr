//go:build integration
// +build integration

package clients

import (
	"context"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestOllamaClient_CreateChatCompletion_Integration(t *testing.T) {
	defaultConfig := openai.DefaultConfig("foobar")
	// default ollama port
	defaultConfig.BaseURL = "http://localhost:11434"
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
		{
			name: "successful",
			fields: fields{
				cfg: defaultConfig,
			},
			args: args{
				ctx: context.Background(),
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
			wantErr: false,
			wantResponse: openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: "foobar",
						},
					},
				},
			},
		},
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
			if len(gotResponse.Choices) == 0 {
				t.Errorf("OllamaClient.CreateChatCompletion() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
