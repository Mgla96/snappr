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
		// TODO: Add test cases.
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
