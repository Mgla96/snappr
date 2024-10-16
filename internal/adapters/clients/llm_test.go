package clients

import (
	"context"
	"errors"
	"testing"

	"github.com/Mgla96/snappr/internal/adapters/clients/clientsfakes"
	snapprErrors "github.com/Mgla96/snappr/internal/errors"
	"github.com/sashabaranov/go-openai"
)

func TestModelToContextWindow(t *testing.T) {
	type args struct {
		model ModelType
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: string(GPT3_5Turbo0125),
			args: args{model: GPT3_5Turbo0125},
			want: 16385,
		},
		{
			name: string(GPT4_turbo),
			args: args{model: GPT4_turbo},
			want: 128000,
		},
		{
			name: "not found",
			args: args{model: "not-found"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ModelToContextWindow(tt.args.model); got != tt.want {
				t.Errorf("ModelToContextWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenAIClient_GenerateChatCompletion(t *testing.T) {
	errCreatingChatCompletion := snapprErrors.New("error creating chat completion")
	type fields struct {
		aiClient openAIBackingClient
	}
	type args struct {
		ctx      context.Context
		messages []openai.ChatCompletionMessage
		model    ModelType
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        string
		wantErr     bool
		wantErrType error
	}{
		{
			name: "error creating chat completion",
			fields: fields{
				aiClient: &clientsfakes.FakeOpenAIBackingClient{
					CreateChatCompletionStub: func(ctx context.Context, request openai.ChatCompletionRequest) (response openai.ChatCompletionResponse, err error) {
						return openai.ChatCompletionResponse{}, errCreatingChatCompletion
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				messages: []openai.ChatCompletionMessage{},
				model:    GPT3_5Turbo0125,
			},
			wantErr:     true,
			wantErrType: errCreatingChatCompletion,
		},
		{
			name: "no chat completion choices returned",
			fields: fields{
				aiClient: &clientsfakes.FakeOpenAIBackingClient{
					CreateChatCompletionStub: func(ctx context.Context, request openai.ChatCompletionRequest) (response openai.ChatCompletionResponse, err error) {
						return openai.ChatCompletionResponse{}, nil
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				messages: []openai.ChatCompletionMessage{},
				model:    GPT3_5Turbo0125,
			},
			wantErr:     true,
			wantErrType: ErrNoChatCompletionChoices,
		},
		{
			name: "success",
			fields: fields{
				aiClient: &clientsfakes.FakeOpenAIBackingClient{
					CreateChatCompletionStub: func(ctx context.Context, request openai.ChatCompletionRequest) (response openai.ChatCompletionResponse, err error) {
						return openai.ChatCompletionResponse{
							Choices: []openai.ChatCompletionChoice{
								{
									Message: openai.ChatCompletionMessage{
										Content: "foobar",
									},
								},
							},
						}, nil
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				messages: []openai.ChatCompletionMessage{},
				model:    GPT3_5Turbo0125,
			},
			want:    "foobar",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oc := &OpenAIClient{
				aiClient: tt.fields.aiClient,
			}
			got, err := oc.GenerateChatCompletion(tt.args.ctx, tt.args.messages, tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenAIClient.GenerateChatCompletion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr && !errors.Is(err, tt.wantErrType) {
				t.Errorf("OpenAIClient.GenerateChatCompletion() error = %v, wantErrType %v", err, tt.wantErrType)
				return
			}
			if got != tt.want {
				t.Errorf("OpenAIClient.GenerateChatCompletion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCustomOpenAIClient(t *testing.T) {
	type args struct {
		authToken string
		baseURL   string
		apiType   APIType
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Valid authToken and baseURL",
			args: args{
				authToken: "foobar",
				baseURL:   "http://foobar.com",
				apiType:   OPENAIAPI,
			},
		},
		{
			name: "Valid authToken and baseURL with OLLAMAAPI",
			args: args{
				authToken: "foobar",
				baseURL:   "http://foobar.com",
				apiType:   OLLAMAAPI,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCustomOpenAIClient(tt.args.authToken, tt.args.baseURL, tt.args.apiType)

			if got.aiClient == nil {
				t.Errorf("NewCustomOpenAIClient() aiClient = nil, want not nil")
			}
		})
	}
}

func TestNewOpenAIClient(t *testing.T) {
	type args struct {
		apiKey string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Valid apiKey",
			args: args{
				apiKey: "foobar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOpenAIClient(tt.args.apiKey)

			if got.aiClient == nil {
				t.Errorf("NewCustomOpenAIClient() aiClient = nil, want not nil")
			}
		})
	}
}
