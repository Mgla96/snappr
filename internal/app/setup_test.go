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
		setupEnvironmentVariables()
		wantCfg := expectedConfig()
		got, err := Setup()
		if err != nil {
			t.Errorf("Setup() error = %v, want %v", err, nil)
		}
		if !reflect.DeepEqual(got.cfg, wantCfg) {
			t.Errorf("Setup() = %v, want %v, diff: %s", got.cfg, wantCfg, cmp.Diff(got.cfg, wantCfg))
		}
		validateClients(got)
	})

	t.Run("test unsuccessful setup", func(t *testing.T) {
		clearEnvironmentVariables()
		got, err := Setup()
		if err == nil {
			t.Error("Expected error but got none")
		}
		if got != nil {
			t.Errorf("Setup() expected configuration to be nil, got: %v", got)
		}
	})
}

func setupEnvironmentVariables() {
	os.Setenv("PR_LOG_LEVEL", "0")
	os.Setenv("PR_GITHUB_TOKEN", "token")
	os.Setenv("PR_GITHUB_OWNER", "owner")
	os.Setenv("PR_GITHUB_REPO", "repo")
	os.Setenv("PR_LLM_TOKEN", "foobar")
	os.Setenv("PR_LLM_ENDPOINT", "http://localhost:8080")
	os.Setenv("PR_INPUT_PROMPT_WORKFLOWS", "foo")
}

func clearEnvironmentVariables() {
	os.Unsetenv("PR_LOG_LEVEL")
	os.Unsetenv("PR_GITHUB_TOKEN")
	os.Unsetenv("PR_GITHUB_OWNER")
	os.Unsetenv("PR_GITHUB_REPO")
	os.Unsetenv("PR_LLM_TOKEN")
	os.Unsetenv("PR_LLM_ENDPOINT")
}

func expectedConfig() *config.Config {
	return &config.Config{
		Log: config.Log{Level: 0},
		Github: config.Github{Token: "token", Owner: "owner", Repo: "repo"},
		LLM: config.LLM{Token: "foobar", Endpoint: "http://localhost:8080", DefaultModel: "gpt-4-turbo", APIType: "openai", Retries: 3},
	}
}

func validateClients(got *ConfigSetup) {
	if got.githubClient == nil {
		t.Errorf("Github client should not be nil")
	}
	if got.llmClient == nil {
		t.Errorf("LLM client should not be nil")
	}
}
