package app

import (
	"reflect"
	"testing"

	"github.com/Mgla96/snappr/internal/config"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		name string
		want *App
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Setup(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Setup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetupNoEnv(t *testing.T) {
	cfg := &config.Config{
		Log: config.Log{
			Level: 0,
		},
		Github: config.Github{
			Token: "token",
			Owner: "owner",
			Repo:  "repo",
		},
		LLM: config.LLM{
			Token:    "foobar",
			Endpoint: "http://localhost:8080",
		},
	}
	cfgNoEndpoint := &config.Config{
		Log: config.Log{
			Level: 0,
		},
		Github: config.Github{
			Token: "token",
			Owner: "owner",
			Repo:  "repo",
		},
		LLM: config.LLM{
			Token: "foobar",
		},
	}
	type args struct {
		cfg *config.Config
	}
	type want struct {
		cfg              *config.Config
		wantGithubClient bool
		wantLlmClient    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test successful setup no env",
			args: args{
				cfg: cfg,
			},
			want: want{
				cfg:              cfg,
				wantGithubClient: true,
				wantLlmClient:    true,
			},
		},
		{
			name: "test successful setup no env and no llm endpoint",
			args: args{
				cfg: cfgNoEndpoint,
			},
			want: want{
				cfg:              cfgNoEndpoint,
				wantGithubClient: true,
				wantLlmClient:    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetupNoEnv(tt.args.cfg)

			if !reflect.DeepEqual(got.cfg, tt.want.cfg) {
				t.Errorf("SetupNoEnv() = %v, want %v", got.cfg, tt.want.cfg)
			}
			if tt.want.wantGithubClient {
				if got.githubClient == nil {
					t.Errorf("SetupNoEnv() = %v, want %v", got.githubClient, tt.want.wantGithubClient)
				}
			}
			if tt.want.wantLlmClient {
				if got.llmClient == nil {
					t.Errorf("SetupNoEnv() = %v, want %v", got.llmClient, tt.want.wantLlmClient)
				}
			}
		})
	}
}
