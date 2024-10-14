package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/Mgla96/snappr/internal/adapters/clients"
	"github.com/rs/zerolog"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name: "valid config",
			envVars: map[string]string{
				"PR_LOG_LEVEL":         "info",
				"PR_GITHUB_TOKEN":      "fake-github-token",
				"PR_GITHUB_OWNER":      "fake-owner",
				"PR_GITHUB_REPO":       "fake-repo",
				"PR_LLM_TOKEN":         "fake-llm-token",
				"PR_LLM_DEFAULT_MODEL": string(clients.GPT4_turbo),
				"PR_LLM_ENDPOINT":      "https://foobar.llm",
			},
			want: &Config{
				Log: Log{
					Level: zerolog.InfoLevel,
				},
				Github: Github{
					Token: "fake-github-token",
					Owner: "fake-owner",
					Repo:  "fake-repo",
				},
				LLM: LLM{
					Token:        "fake-llm-token",
					DefaultModel: clients.GPT4_turbo,
					Endpoint:     "https://foobar.llm",
					APIType:      clients.OPENAIAPI,
				},
			},
			wantErr: false,
		},
		{
			name: "missing required env vars",
			envVars: map[string]string{
				"PR_LOG_LEVEL": "info",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			got, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
