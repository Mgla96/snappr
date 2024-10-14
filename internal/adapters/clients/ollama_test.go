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
