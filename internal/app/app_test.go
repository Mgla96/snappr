package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/Mgla96/snappr/internal/adapters/clients"
	"github.com/Mgla96/snappr/internal/app/appfakes"
	"github.com/Mgla96/snappr/internal/config"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v39/github"
	"github.com/rs/zerolog"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/require"
)

func TestApp_ExecuteCreatePR(t *testing.T) {
	examplePRCreation := PRCreation{
		Title: "foo",
		Body:  "bar",
		UpdatedFiles: []PRCreationFile{
			{
				Path:          "foo",
				FullContent:   "bar",
				CommitMessage: "baz",
			},
		},
	}
	prCreationBytes, err := json.Marshal(examplePRCreation)
	if err != nil {
		t.Fatalf("failed to marshal examplePRCreation: %v", err)
	}
	type fields struct {
		cfg          *config.Config
		githubClient githubClient
		llmClient    llmClient
		log          zerolog.Logger
	}
	type args struct {
		ctx          context.Context
		commitSHA    string
		branch       string
		workflowName string
		printOnly    bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Error getting commit code",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetCommitCodeStub: func(context.Context, string, string, string, clients.CodeFilter) (map[string]string, error) {
						return map[string]string{}, fmt.Errorf("error getting commit code")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Empty commit code",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetCommitCodeStub: func(context.Context, string, string, string, clients.CodeFilter) (map[string]string, error) {
						return map[string]string{}, nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error generating chat completion",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetCommitCodeStub: func(context.Context, string, string, string, clients.CodeFilter) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return "", fmt.Errorf("error generating chat completion")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Empty response from LLM",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetCommitCodeStub: func(context.Context, string, string, string, clients.CodeFilter) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return "", nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Print only",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetCommitCodeStub: func(context.Context, string, string, string, clients.CodeFilter) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return string(prCreationBytes), nil
					},
				},
			},
			args: args{
				printOnly: true,
			},
			wantErr: false,
		},
		{
			name: "Error unmarshal LLM Response to PrCreation struct",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetCommitCodeStub: func(context.Context, string, string, string, clients.CodeFilter) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return "{'foo': 'bar'}", nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error creating Github Branch",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetCommitCodeStub: func(context.Context, string, string, string, clients.CodeFilter) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
					CreateBranchStub: func(context.Context, string, string, string, string) error {
						return fmt.Errorf("error creating branch")
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return string(prCreationBytes), nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error adding commit to branch",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetCommitCodeStub: func(context.Context, string, string, string, clients.CodeFilter) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
					CreateBranchStub: func(context.Context, string, string, string, string) error {
						return nil
					},
					AddCommitToBranchStub: func(context.Context, string, string, string, string, string, []byte) error {
						return fmt.Errorf("error adding commit to branch")
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return string(prCreationBytes), nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error creating pull request",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetCommitCodeStub: func(context.Context, string, string, string, clients.CodeFilter) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
					CreateBranchStub: func(context.Context, string, string, string, string) error {
						return nil
					},
					AddCommitToBranchStub: func(context.Context, string, string, string, string, string, []byte) error {
						return nil
					},
					CreatePullRequestStub: func(context.Context, string, string, string, string, string, string) (*github.PullRequest, error) {
						return nil, fmt.Errorf("error creating pull request")
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return string(prCreationBytes), nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "created pull request",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetCommitCodeStub: func(context.Context, string, string, string, clients.CodeFilter) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
					CreateBranchStub: func(context.Context, string, string, string, string) error {
						return nil
					},
					AddCommitToBranchStub: func(context.Context, string, string, string, string, string, []byte) error {
						return nil
					},
					CreatePullRequestStub: func(context.Context, string, string, string, string, string, string) (*github.PullRequest, error) {
						return &github.PullRequest{
							HTMLURL: github.String(""),
						}, nil
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return string(prCreationBytes), nil
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				cfg:          tt.fields.cfg,
				githubClient: tt.fields.githubClient,
				llmClient:    tt.fields.llmClient,
				log:          tt.fields.log,
			}
			if err := a.ExecuteCreatePR(tt.args.ctx, tt.args.commitSHA, tt.args.branch, tt.args.workflowName, `.*\.go$`, tt.args.printOnly); (err != nil) != tt.wantErr {
				t.Errorf("App.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseRawDiff(t *testing.T) {
	rawDiff := `diff --git a/file1.go b/file1.go
index abc123..def456 100644
--- a/file1.go
+++ b/file1.go
@@ -1,4 +1,6 @@
 line1
 line2
+added line3
 line4
+added line5

diff --git a/file2.go b/file2.go
index 789abc..012def 100644
--- a/file2.go
+++ b/file2.go
@@ -10,7 +10,8 @@
 func exampleFunction() {
     line1
-    removed line2
+    added line2
     line3
 }`
	type args struct {
		diff string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "parse raw diff",
			args: args{
				diff: rawDiff,
			},
			want: map[string]string{
				"file1.go": `index abc123..def456 100644
--- a/file1.go
+++ b/file1.go
@@ -1,4 +1,6 @@
 line1
 line2
+added line3
 line4
+added line5

`,
				"file2.go": `index 789abc..012def 100644
--- a/file2.go
+++ b/file2.go
@@ -10,7 +10,8 @@
 func exampleFunction() {
     line1
-    removed line2
+    added line2
     line3
 }`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRawDiff(tt.args.diff); !cmp.Equal(got, tt.want) {
				t.Errorf("parseRawDiff() diff = %s", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestApp_ExecutePRReview(t *testing.T) {
	rawDiff := `diff --git a/file1.go b/file1.go
index abc123..def456 100644
--- a/file1.go
+++ b/file1.go
@@ -1,4 +1,6 @@
 line1
 line2
+added line3
 line4
+added line5

diff --git a/file2.go b/file2.go
index 789abc..012def 100644
--- a/file2.go
+++ b/file2.go
@@ -10,7 +10,8 @@
 func exampleFunction() {
     line1
-    removed line2
+    added line2
     line3
 }`

	examplePRReview := PRReviewMap{
		"file1.go": []PRCommentInfo{
			{
				CommentBody: "foo",
				StartLine:   1,
				Line:        2,
			},
		},
	}
	examplePRReviewBytes, err := json.Marshal(examplePRReview)
	if err != nil {
		t.Fatalf("failed to marshal examplePRReview: %v", err)
	}

	type fields struct {
		cfg          *config.Config
		githubClient githubClient
		llmClient    llmClient
		log          zerolog.Logger
	}
	type args struct {
		ctx          context.Context
		commitSHA    string
		prNumber     int
		workflowName string
		printOnly    bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Error getting PR code",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetPRCodeStub: func(context.Context, string, string, int, *github.ListOptions) (map[string]string, error) {
						return map[string]string{}, fmt.Errorf("error getting PR code")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Empty PR code",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetPRCodeStub: func(context.Context, string, string, int, *github.ListOptions) (map[string]string, error) {
						return map[string]string{}, nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error getting PR Patch",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetPRCodeStub: func(context.Context, string, string, int, *github.ListOptions) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
					GetPRPatchStub: func(context.Context, string, string, int) (string, error) {
						return "", fmt.Errorf("error getting PR patch")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Empty prompt workflow and failure to generate chat completion",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetPRCodeStub: func(context.Context, string, string, int, *github.ListOptions) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
					GetPRPatchStub: func(context.Context, string, string, int) (string, error) {
						return rawDiff, nil
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return "", fmt.Errorf("error generating chat completion")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Empty response from GenerateChatCompletion",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetPRCodeStub: func(context.Context, string, string, int, *github.ListOptions) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
					GetPRPatchStub: func(context.Context, string, string, int) (string, error) {
						return rawDiff, nil
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return "", nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "print only",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetPRCodeStub: func(context.Context, string, string, int, *github.ListOptions) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
					GetPRPatchStub: func(context.Context, string, string, int) (string, error) {
						return rawDiff, nil
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return string(examplePRReviewBytes), nil
					},
				},
			},
			args: args{
				printOnly: true,
			},
			wantErr: false,
		},
		{
			name: "error unmarshalling PR Review Map",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetPRCodeStub: func(context.Context, string, string, int, *github.ListOptions) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
					GetPRPatchStub: func(context.Context, string, string, int) (string, error) {
						return rawDiff, nil
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return strings.TrimRight(string(examplePRReviewBytes), "}"), nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error adding comment to pull request review",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetPRCodeStub: func(context.Context, string, string, int, *github.ListOptions) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
					GetPRPatchStub: func(context.Context, string, string, int) (string, error) {
						return rawDiff, nil
					},
					AddCommentToPullRequestReviewStub: func(context.Context, string, string, int, string, string, string, int, int) (*github.PullRequestComment, error) {
						return nil, fmt.Errorf("error adding comment to pull request review")
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return string(examplePRReviewBytes), nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "successful PR review",
			fields: fields{
				cfg: &config.Config{},
				githubClient: &appfakes.FakeGithubClient{
					GetPRCodeStub: func(context.Context, string, string, int, *github.ListOptions) (map[string]string, error) {
						return map[string]string{
							"file1.go": "code1",
						}, nil
					},
					GetPRPatchStub: func(context.Context, string, string, int) (string, error) {
						return rawDiff, nil
					},
					AddCommentToPullRequestReviewStub: func(context.Context, string, string, int, string, string, string, int, int) (*github.PullRequestComment, error) {
						return &github.PullRequestComment{
							Body: github.String(""),
						}, nil
					},
				},
				llmClient: &appfakes.FakeLlmClient{
					GenerateChatCompletionStub: func(context.Context, []openai.ChatCompletionMessage, clients.ModelType) (string, error) {
						return string(examplePRReviewBytes), nil
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				cfg:          tt.fields.cfg,
				githubClient: tt.fields.githubClient,
				llmClient:    tt.fields.llmClient,
				log:          tt.fields.log,
			}
			if err := a.ExecutePRReview(tt.args.ctx, tt.args.commitSHA, tt.args.prNumber, tt.args.workflowName, tt.args.printOnly); (err != nil) != tt.wantErr {
				t.Errorf("App.ExecutePRReview() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckGHToken(t *testing.T) {
	t.Run("Missing GH_TOKEN", func(t *testing.T) {
		os.Unsetenv("GH_TOKEN")
		err := CheckGHToken()
		require.Error(t, err)
		require.EqualError(t, err, "GH_TOKEN environment variable is required. Please set it before running this command")
	})
	t.Run("GH_TOKEN set", func(t *testing.T) {
		os.Setenv("GH_TOKEN", "foobar")
		err := CheckGHToken()
		require.NoError(t, err)
	})
}

func TestCheckLLMToken(t *testing.T) {
	t.Run("Missing LLM_TOKEN", func(t *testing.T) {
		os.Unsetenv("LLM_TOKEN")
		err := CheckLLMToken()
		require.Error(t, err)
		require.EqualError(t, err, "LLM_TOKEN environment variable is required. Please set it before running this command")
	})
	t.Run("LLM_TOKEN set", func(t *testing.T) {
		os.Setenv("LLM_TOKEN", "foobar")
		err := CheckLLMToken()
		require.NoError(t, err)
	})
}

func TestCheckTokens(t *testing.T) {
	t.Run("Missing GH_TOKEN", func(t *testing.T) {
		os.Unsetenv("GH_TOKEN")
		os.Setenv("LLM_TOKEN", "foobar")
		err := CheckTokens()
		require.Error(t, err)
		require.EqualError(t, err, "error no github token: GH_TOKEN environment variable is required. Please set it before running this command")
	})
	t.Run("Missing LLM_TOKEN", func(t *testing.T) {
		os.Setenv("GH_TOKEN", "foobar")
		os.Unsetenv("LLM_TOKEN")
		err := CheckTokens()
		require.Error(t, err)
		require.EqualError(t, err, "error no llm token: LLM_TOKEN environment variable is required. Please set it before running this command")
	})

	t.Run("GH_TOKEN and LLM token set", func(t *testing.T) {
		os.Setenv("GH_TOKEN", "foobar")
		os.Setenv("LLM_TOKEN", "foobar")
		err := CheckTokens()
		require.NoError(t, err)
	})
}

func TestNew(t *testing.T) {
	cfg := &config.Config{
		Log: config.Log{
			Level: zerolog.InfoLevel,
		},
		Github: config.Github{
			Token: "foobar",
			Owner: "foo",
			Repo:  "bar",
		},
		LLM: config.LLM{
			Token: "foobar",
		},
	}
	ghc := &appfakes.FakeGithubClient{}
	llmc := &appfakes.FakeLlmClient{}
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	type args struct {
		cfg          *config.Config
		githubClient githubClient
		llmClient    llmClient
		logger       zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    *App
		wantErr bool
	}{
		{
			name: "New App",
			args: args{
				cfg:          cfg,
				githubClient: ghc,
				llmClient:    llmc,
				logger:       logger,
			},
			want: &App{
				cfg:          cfg,
				githubClient: ghc,
				llmClient:    llmc,
				log:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.cfg, tt.args.githubClient, tt.args.llmClient, tt.args.logger)
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
