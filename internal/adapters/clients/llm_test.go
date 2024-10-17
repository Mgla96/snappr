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
			name: "GPT3.5 Turbo", args: args{model: GPT3_5Turbo0125}, want: 16384,
			name: "GPT4 Turbo", args: args{model: GPT4_turbo}, want: 128000,
			name: "not found", args: args{model: "not-found"}, want: -1,
		},
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := ModelToContextWindow(tt.args.model)
				if got != tt.want {
					t.Errorf("ModelToContextWindow() = %d, want %d", got, tt.want)
				}
			})
		}
	}

// Simplify and correct API client usage scenarios
func TestOpenAIClient_GenerateChatCompletion(t *testing.T) {...}
func TestNewCustomOpenAIClient(t *testing.T) {...}
func TestNewOpenAIClient(t *testing.T) {...}
