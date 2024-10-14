package app

import (
	"os"
	"reflect"
	"testing"

	"github.com/Mgla96/snappr/internal/config"
	"github.com/google/go-cmp/cmp"
)

func TestSetup(t *testing.T) {

	t.Run("test successful setup", func(t *testing.T) {
		wantCfg := &config.Config{
			Log: config.Log{
				Level: 0,
			},
			Github: config.Github{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
			},
			LLM: config.LLM{
				Token:        "foobar",
				Endpoint:     "http://localhost:8080",
				DefaultModel: "gpt-4-turbo",
				APIType:      "openai",
			},
		}

		os.Setenv("PR_LOG_LEVEL", "0")
		os.Setenv("PR_GITHUB_TOKEN", "token")
		os.Setenv("PR_GITHUB_OWNER", "owner")
		os.Setenv("PR_GITHUB_REPO", "repo")
		os.Setenv("PR_LLM_TOKEN", "foobar")
		os.Setenv("PR_LLM_ENDPOINT", "http://localhost:8080")
		os.Setenv("PR_INPUT_PROMPT_WORKFLOWS", "foo")

		got, err := Setup()
		if err != nil {
			t.Errorf("Setup() error = %v, want %v", err, nil)
		}

		if !reflect.DeepEqual(got.cfg, wantCfg) {
			t.Errorf("SetupNoEnv() = %v, want %v, diff: %s", got.cfg, wantCfg, cmp.Diff(got.cfg, wantCfg))
		}
		if got.githubClient == nil {
			t.Errorf("SetupNoEnv() = %v, want %v", got.githubClient, true)
		}
		if got.llmClient == nil {
			t.Errorf("SetupNoEnv() = %v, want %v", got.llmClient, true)
		}

	})

	t.Run("test unsuccessful setup", func(t *testing.T) {
		os.Unsetenv("PR_LOG_LEVEL")
		os.Unsetenv("PR_GITHUB_TOKEN")
		os.Unsetenv("PR_GITHUB_OWNER")
		os.Unsetenv("PR_GITHUB_REPO")
		os.Unsetenv("PR_LLM_TOKEN")
		os.Unsetenv("PR_LLM_ENDPOINT")

		got, err := Setup()
		if err == nil {
			t.Error("Setup() error nil but wanted error", err)
		}
		if got != nil {
			t.Errorf("Setup() = %v, want %v", got.cfg, nil)
		}
	})

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
