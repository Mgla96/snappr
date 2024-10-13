package clients

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/Mgla96/snappr/internal/adapters/clients/clientsfakes"
	"github.com/google/go-github/v39/github"
	"github.com/rs/zerolog"
)

func TestIsDoNotEditFile(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Contains DO NOT EDIT",
			args: args{
				data: []byte("DO NOT EDIT\nFoobar"),
			},
			want: true,
		},
		{
			name: "Empty file",
			args: args{
				data: []byte(""),
			},
			want: false,
		},
		{
			name: "Partial match",
			args: args{
				data: []byte("DO NO EDIT\nSome other content"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDoNotEditFile(tt.args.data); got != tt.want {
				t.Errorf("IsDoNotEditFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubClient_processEntry(t *testing.T) {
	type fields struct {
		ghGitClient         gitService
		ghPullRequestClient pullRequestService
		log                 zerolog.Logger
	}
	type args struct {
		entry      *github.TreeEntry
		codeFilter CodeFilter
		context    context.Context
		owner      string
		repo       string
		files      map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "entry type not blob",
			args: args{
				entry: &github.TreeEntry{
					Type: github.String("tree"),
				},
			},
			wantErr: false,
		},
		{
			name: "regex path not matching",
			args: args{
				entry: &github.TreeEntry{
					Type: github.String("blob"),
					Path: github.String("file.txt"),
				},
				codeFilter: CodeFilter{
					FileRegexPattern: ".*\\.go$",
				},
			},
			wantErr: false,
		},
		{
			name: "regex path matching and error getting blob",
			fields: fields{
				ghGitClient: &clientsfakes.FakeGitService{
					GetBlobStub: func(context.Context, string, string, string) (*github.Blob, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			args: args{
				entry: &github.TreeEntry{
					Type: github.String("blob"),
					Path: github.String("file.go"),
				},
				codeFilter: CodeFilter{
					FileRegexPattern: ".*\\.go$",
				},
			},
			wantErr: true,
		},
		// {
		// 	name: "get blob unexpected status",
		// 	fields: fields{
		// 		ghGitClient: &clientsfakes.FakeGitService{
		// 			GetBlobStub: func(context.Context, string, string, string) (*github.Blob, *github.Response, error) {
		// 				resp := &github.Response{}
		// 				resp.StatusCode = 500
		// 				return nil, resp, fmt.Errorf("error")
		// 			},
		// 		},
		// 	},
		// 	args: args{
		// 		entry: &github.TreeEntry{
		// 			Type: github.String("blob"),
		// 			Path: github.String("file.go"),
		// 		},
		// 		codeFilter: CodeFilter{
		// 			FileRegexPattern: ".*\\.go$",
		// 		},
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GithubClient{
				ghGitClient:         tt.fields.ghGitClient,
				ghPullRequestClient: tt.fields.ghPullRequestClient,
				log:                 tt.fields.log,
			}
			if err := gc.processEntry(tt.args.entry, tt.args.codeFilter, tt.args.context, tt.args.owner, tt.args.repo, tt.args.files); (err != nil) != tt.wantErr {
				t.Errorf("GithubClient.processEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGithubClient_CreateBranch(t *testing.T) {
	type fields struct {
		ghGitClient         gitService
		ghPullRequestClient pullRequestService
		log                 zerolog.Logger
	}
	type args struct {
		ctx        context.Context
		owner      string
		repo       string
		newBranch  string
		baseBranch string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error getting reference",
			fields: fields{
				ghGitClient: &clientsfakes.FakeGitService{
					GetRefStub: func(context.Context, string, string, string) (*github.Reference, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			args: args{
				ctx:       context.Background(),
				owner:     "owner",
				repo:      "repo",
				newBranch: "newBranch",
			},
			wantErr: true,
		},
		{
			name: "error creating reference",
			fields: fields{
				ghGitClient: &clientsfakes.FakeGitService{
					GetRefStub: func(context.Context, string, string, string) (*github.Reference, *github.Response, error) {
						return &github.Reference{
							Object: &github.GitObject{
								SHA: github.String("sha"),
							},
						}, nil, nil
					},
					CreateRefStub: func(context.Context, string, string, *github.Reference) (*github.Reference, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			args: args{
				ctx:       context.Background(),
				owner:     "owner",
				repo:      "repo",
				newBranch: "newBranch",
			},
			wantErr: true,
		},
		{
			name: "successful branch creation",
			fields: fields{
				ghGitClient: &clientsfakes.FakeGitService{
					GetRefStub: func(context.Context, string, string, string) (*github.Reference, *github.Response, error) {
						return &github.Reference{
							Object: &github.GitObject{
								SHA: github.String("sha"),
							},
						}, nil, nil
					},
					CreateRefStub: func(context.Context, string, string, *github.Reference) (*github.Reference, *github.Response, error) {
						return nil, nil, nil
					},
				},
			},
			args: args{
				ctx:       context.Background(),
				owner:     "owner",
				repo:      "repo",
				newBranch: "newBranch",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GithubClient{
				ghGitClient:         tt.fields.ghGitClient,
				ghPullRequestClient: tt.fields.ghPullRequestClient,
				log:                 tt.fields.log,
			}
			if err := gc.CreateBranch(tt.args.ctx, tt.args.owner, tt.args.repo, tt.args.newBranch, tt.args.baseBranch); (err != nil) != tt.wantErr {
				t.Errorf("GithubClient.CreateBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGithubClient_GetLatestCommitFromBranch(t *testing.T) {
	type fields struct {
		ghGitClient         gitService
		ghPullRequestClient pullRequestService
		log                 zerolog.Logger
	}
	type args struct {
		ctx    context.Context
		owner  string
		repo   string
		branch string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "error getting reference",
			fields: fields{
				ghGitClient: &clientsfakes.FakeGitService{
					GetRefStub: func(context.Context, string, string, string) (*github.Reference, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			args: args{
				ctx:    context.Background(),
				owner:  "owner",
				repo:   "repo",
				branch: "branch",
			},
			wantErr: true,
		},
		{
			name: "successful branch creation",
			fields: fields{
				ghGitClient: &clientsfakes.FakeGitService{
					GetRefStub: func(context.Context, string, string, string) (*github.Reference, *github.Response, error) {
						return &github.Reference{
							Object: &github.GitObject{
								SHA: github.String("sha"),
							},
						}, nil, nil
					},
				},
			},
			args: args{
				ctx:    context.Background(),
				owner:  "owner",
				repo:   "repo",
				branch: "branch",
			},
			wantErr: false,
			want:    "sha",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GithubClient{
				ghGitClient:         tt.fields.ghGitClient,
				ghPullRequestClient: tt.fields.ghPullRequestClient,
				log:                 tt.fields.log,
			}
			got, err := gc.GetLatestCommitFromBranch(tt.args.ctx, tt.args.owner, tt.args.repo, tt.args.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubClient.GetLatestCommitFromBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GithubClient.GetLatestCommitFromBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubClient_AddCommitToBranch(t *testing.T) {
	type fields struct {
		ghGitClient         gitService
		ghPullRequestClient pullRequestService
		log                 zerolog.Logger
	}
	type args struct {
		ctx           context.Context
		owner         string
		repo          string
		branch        string
		filePath      string
		commitMessage string
		fileContent   []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error getting reference",
			fields: fields{
				ghGitClient: &clientsfakes.FakeGitService{
					GetRefStub: func(context.Context, string, string, string) (*github.Reference, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			args: args{
				ctx:    context.Background(),
				owner:  "owner",
				repo:   "repo",
				branch: "branch",
			},
			wantErr: true,
		},
		{
			name: "error creating blob",
			fields: fields{
				ghGitClient: &clientsfakes.FakeGitService{
					GetRefStub: func(context.Context, string, string, string) (*github.Reference, *github.Response, error) {
						return &github.Reference{
							Object: &github.GitObject{
								SHA: github.String("sha"),
							},
						}, nil, nil
					},
					CreateBlobStub: func(context.Context, string, string, *github.Blob) (*github.Blob, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			args: args{
				ctx:    context.Background(),
				owner:  "owner",
				repo:   "repo",
				branch: "branch",
			},
			wantErr: true,
		},
		{
			name: "error getting tree",
			fields: fields{
				ghGitClient: &clientsfakes.FakeGitService{
					GetRefStub: func(context.Context, string, string, string) (*github.Reference, *github.Response, error) {
						return &github.Reference{
							Object: &github.GitObject{
								SHA: github.String("sha"),
							},
						}, nil, nil
					},
					CreateBlobStub: func(context.Context, string, string, *github.Blob) (*github.Blob, *github.Response, error) {
						return &github.Blob{
							SHA: github.String("blobSHA"),
						}, nil, nil
					},
					GetTreeStub: func(context.Context, string, string, string, bool) (*github.Tree, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				owner:    "owner",
				repo:     "repo",
				branch:   "branch",
				filePath: "filePath",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GithubClient{
				ghGitClient:         tt.fields.ghGitClient,
				ghPullRequestClient: tt.fields.ghPullRequestClient,
				log:                 tt.fields.log,
			}
			if err := gc.AddCommitToBranch(tt.args.ctx, tt.args.owner, tt.args.repo, tt.args.branch, tt.args.filePath, tt.args.commitMessage, tt.args.fileContent); (err != nil) != tt.wantErr {
				t.Errorf("GithubClient.AddCommitToBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGithubClient_GetCommitCode(t *testing.T) {
	type fields struct {
		ghGitClient         gitService
		ghPullRequestClient pullRequestService
		log                 zerolog.Logger
	}
	type args struct {
		context    context.Context
		owner      string
		repo       string
		commitSHA  string
		codeFilter CodeFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "error getting commit",
			fields: fields{
				ghGitClient: &clientsfakes.FakeGitService{
					GetCommitStub: func(context.Context, string, string, string) (*github.Commit, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error getting tree",
			fields: fields{
				ghGitClient: &clientsfakes.FakeGitService{
					GetCommitStub: func(context.Context, string, string, string) (*github.Commit, *github.Response, error) {
						return &github.Commit{
							SHA: github.String("sha"),
						}, nil, nil
					},
					GetTreeStub: func(context.Context, string, string, string, bool) (*github.Tree, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GithubClient{
				ghGitClient:         tt.fields.ghGitClient,
				ghPullRequestClient: tt.fields.ghPullRequestClient,
				log:                 tt.fields.log,
			}
			got, err := gc.GetCommitCode(tt.args.context, tt.args.owner, tt.args.repo, tt.args.commitSHA, tt.args.codeFilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubClient.GetCommitCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GithubClient.GetCommitCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubClient_AddCommentToPullRequestReview(t *testing.T) {
	type fields struct {
		ghGitClient         gitService
		ghPullRequestClient pullRequestService
		log                 zerolog.Logger
	}
	type args struct {
		ctx         context.Context
		owner       string
		repo        string
		prNumber    int
		commentBody string
		commitID    string
		path        string
		startLine   int
		line        int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *github.PullRequestComment
		wantErr bool
	}{
		{
			name: "error creating comment",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					CreateCommentStub: func(context.Context, string, string, int, *github.PullRequestComment) (*github.PullRequestComment, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "successfully creating comment",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					CreateCommentStub: func(context.Context, string, string, int, *github.PullRequestComment) (*github.PullRequestComment, *github.Response, error) {
						return &github.PullRequestComment{
							Body: github.String("commentBody"),
						}, nil, nil
					},
				},
			},
			want:    &github.PullRequestComment{Body: github.String("commentBody")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GithubClient{
				ghGitClient:         tt.fields.ghGitClient,
				ghPullRequestClient: tt.fields.ghPullRequestClient,
				log:                 tt.fields.log,
			}
			got, err := gc.AddCommentToPullRequestReview(tt.args.ctx, tt.args.owner, tt.args.repo, tt.args.prNumber, tt.args.commentBody, tt.args.commitID, tt.args.path, tt.args.startLine, tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubClient.AddCommentToPullRequestReview() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GithubClient.AddCommentToPullRequestReview() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubClient_CreatePullRequest(t *testing.T) {
	type fields struct {
		ghGitClient         gitService
		ghPullRequestClient pullRequestService
		log                 zerolog.Logger
	}
	type args struct {
		ctx   context.Context
		owner string
		repo  string
		title string
		head  string
		base  string
		body  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *github.PullRequest
		wantErr bool
	}{
		{
			name: "error creating pull request",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					CreateStub: func(context.Context, string, string, *github.NewPullRequest) (*github.PullRequest, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			args: args{
				ctx:   context.Background(),
				title: "title",
				head:  "head",
				base:  "base",
				body:  "body",
			},
			wantErr: true,
		},
		{
			name: "successfully create pull request",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					CreateStub: func(context.Context, string, string, *github.NewPullRequest) (*github.PullRequest, *github.Response, error) {
						return &github.PullRequest{
							Title: github.String("title"),
						}, nil, nil
					},
				},
			},
			args: args{
				ctx:   context.Background(),
				title: "title",
				head:  "head",
				base:  "base",
				body:  "body",
			},
			wantErr: false,
			want:    &github.PullRequest{Title: github.String("title")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GithubClient{
				ghGitClient:         tt.fields.ghGitClient,
				ghPullRequestClient: tt.fields.ghPullRequestClient,
				log:                 tt.fields.log,
			}
			got, err := gc.CreatePullRequest(tt.args.ctx, tt.args.owner, tt.args.repo, tt.args.title, tt.args.head, tt.args.base, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubClient.CreatePullRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GithubClient.CreatePullRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubClient_MergePullRequest(t *testing.T) {
	type fields struct {
		ghGitClient         gitService
		ghPullRequestClient pullRequestService
		log                 zerolog.Logger
	}
	type args struct {
		ctx           context.Context
		owner         string
		repo          string
		prNumber      int
		commitMessage string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *github.PullRequestMergeResult
		wantErr bool
	}{
		{
			name: "error merging pull request",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					MergeStub: func(context.Context, string, string, int, string, *github.PullRequestOptions) (*github.PullRequestMergeResult, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "successfully merging pull request",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					MergeStub: func(context.Context, string, string, int, string, *github.PullRequestOptions) (*github.PullRequestMergeResult, *github.Response, error) {
						return &github.PullRequestMergeResult{
							Merged: github.Bool(true),
						}, nil, nil
					},
				},
			},
			wantErr: false,
			want:    &github.PullRequestMergeResult{Merged: github.Bool(true)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GithubClient{
				ghGitClient:         tt.fields.ghGitClient,
				ghPullRequestClient: tt.fields.ghPullRequestClient,
				log:                 tt.fields.log,
			}
			got, err := gc.MergePullRequest(tt.args.ctx, tt.args.owner, tt.args.repo, tt.args.prNumber, tt.args.commitMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubClient.MergePullRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GithubClient.MergePullRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubClient_ListPullRequests(t *testing.T) {
	type fields struct {
		ghGitClient         gitService
		ghPullRequestClient pullRequestService
		log                 zerolog.Logger
	}
	type args struct {
		ctx   context.Context
		owner string
		repo  string
		opts  *github.PullRequestListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*github.PullRequest
		wantErr bool
	}{
		{
			name: "error listing pull requests",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					ListStub: func(context.Context, string, string, *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "successfully listing pull requests",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					ListStub: func(context.Context, string, string, *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
						return []*github.PullRequest{
							{
								Title: github.String("title"),
							},
						}, nil, nil
					},
				},
			},
			wantErr: false,
			want: []*github.PullRequest{
				{
					Title: github.String("title"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GithubClient{
				ghGitClient:         tt.fields.ghGitClient,
				ghPullRequestClient: tt.fields.ghPullRequestClient,
				log:                 tt.fields.log,
			}
			got, err := gc.ListPullRequests(tt.args.ctx, tt.args.owner, tt.args.repo, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubClient.ListPullRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GithubClient.ListPullRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubClient_GetPRCode(t *testing.T) {
	type fields struct {
		ghGitClient         gitService
		ghPullRequestClient pullRequestService
		log                 zerolog.Logger
	}
	type args struct {
		ctx      context.Context
		owner    string
		repo     string
		prNumber int
		opts     *github.ListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "error listing pull request files",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					ListFilesStub: func(context.Context, string, string, int, *github.ListOptions) ([]*github.CommitFile, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no pull request files",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					ListFilesStub: func(context.Context, string, string, int, *github.ListOptions) ([]*github.CommitFile, *github.Response, error) {
						return []*github.CommitFile{
							{
								Filename: github.String("foo.txt"),
							},
						}, nil, nil
					},
				},
			},
			wantErr: false,
			want:    map[string]string{},
		},
		{
			name: "error getting blob",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					ListFilesStub: func(context.Context, string, string, int, *github.ListOptions) ([]*github.CommitFile, *github.Response, error) {
						return []*github.CommitFile{
							{
								Filename: github.String("foo.go"),
								SHA:      github.String("sha"),
							},
						}, nil, nil
					},
				},
				ghGitClient: &clientsfakes.FakeGitService{
					GetBlobStub: func(context.Context, string, string, string) (*github.Blob, *github.Response, error) {
						return nil, nil, fmt.Errorf("error")
					},
				},
			},
			wantErr: true,
		},
		// {
		// 	name: "unexpected status code getting blob",
		// 	fields: fields{
		// 		ghPullRequestClient: &clientsfakes.FakePullRequestService{
		// 			ListFilesStub: func(context.Context, string, string, int, *github.ListOptions) ([]*github.CommitFile, *github.Response, error) {
		// 				return []*github.CommitFile{
		// 					{
		// 						Filename: github.String("foo.go"),
		// 						SHA:      github.String("0d6a88a33a574143c94090cb211e410e6b091d4b"),
		// 					},
		// 				}, nil, nil
		// 			},
		// 		},
		// 		ghGitClient: &clientsfakes.FakeGitService{
		// 			GetBlobStub: func(context.Context, string, string, string) (*github.Blob, *github.Response, error) {
		// 				resp := &github.Response{}
		// 				resp.StatusCode = 500
		// 				return nil, resp, nil
		// 			},
		// 		},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "successful getting blob",
		// 	fields: fields{
		// 		ghPullRequestClient: &clientsfakes.FakePullRequestService{
		// 			ListFilesStub: func(context.Context, string, string, int, *github.ListOptions) ([]*github.CommitFile, *github.Response, error) {
		// 				return []*github.CommitFile{
		// 					{
		// 						Filename: github.String("foo.go"),
		// 						SHA:      github.String("sha"),
		// 					},
		// 				}, nil, nil
		// 			},
		// 		},
		// 		ghGitClient: &clientsfakes.FakeGitService{
		// 			GetBlobStub: func(context.Context, string, string, string) (*github.Blob, *github.Response, error) {
		// 				return &github.Blob{
		// 					Content: github.String("content"),
		// 					SHA:     github.String("sha"),
		// 				}, nil, nil
		// 			},
		// 		},
		// 	},
		// 	wantErr: false,
		// 	want: map[string]string{
		// 		"foo.go": "content",
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GithubClient{
				ghGitClient:         tt.fields.ghGitClient,
				ghPullRequestClient: tt.fields.ghPullRequestClient,
				log:                 tt.fields.log,
			}
			got, err := gc.GetPRCode(tt.args.ctx, tt.args.owner, tt.args.repo, tt.args.prNumber, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubClient.GetPRCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GithubClient.GetPRCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubClient_GetPRDiff(t *testing.T) {
	type fields struct {
		ghGitClient         gitService
		ghPullRequestClient pullRequestService
		log                 zerolog.Logger
	}
	type args struct {
		ctx      context.Context
		owner    string
		repo     string
		prNumber int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "error getting raw",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					GetRawStub: func(context.Context, string, string, int, github.RawOptions) (string, *github.Response, error) {
						return "", nil, fmt.Errorf("error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "successful getting raw",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					GetRawStub: func(context.Context, string, string, int, github.RawOptions) (string, *github.Response, error) {
						return "diff", nil, nil
					},
				},
			},
			want: "diff",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GithubClient{
				ghGitClient:         tt.fields.ghGitClient,
				ghPullRequestClient: tt.fields.ghPullRequestClient,
				log:                 tt.fields.log,
			}
			got, err := gc.GetPRDiff(tt.args.ctx, tt.args.owner, tt.args.repo, tt.args.prNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubClient.GetPRDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GithubClient.GetPRDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubClient_GetPRPatch(t *testing.T) {
	type fields struct {
		ghGitClient         gitService
		ghPullRequestClient pullRequestService
		log                 zerolog.Logger
	}
	type args struct {
		ctx      context.Context
		owner    string
		repo     string
		prNumber int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "error getting raw",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					GetRawStub: func(context.Context, string, string, int, github.RawOptions) (string, *github.Response, error) {
						return "", nil, fmt.Errorf("error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "successful getting raw",
			fields: fields{
				ghPullRequestClient: &clientsfakes.FakePullRequestService{
					GetRawStub: func(context.Context, string, string, int, github.RawOptions) (string, *github.Response, error) {
						return "patch", nil, nil
					},
				},
			},
			wantErr: false,
			want:    "patch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GithubClient{
				ghGitClient:         tt.fields.ghGitClient,
				ghPullRequestClient: tt.fields.ghPullRequestClient,
				log:                 tt.fields.log,
			}
			got, err := gc.GetPRPatch(tt.args.ctx, tt.args.owner, tt.args.repo, tt.args.prNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubClient.GetPRPatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GithubClient.GetPRPatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuffle_IntSlice(t *testing.T) {
	original := []int{1, 2, 3, 4, 5}
	shuffled := make([]int, len(original))
	copy(shuffled, original)

	Shuffle(shuffled)

	if len(original) != len(shuffled) {
		t.Errorf("Expected the shuffled slice length to be %d, got %d", len(original), len(shuffled))
	}

	if reflect.DeepEqual(original, shuffled) {
		t.Errorf("The slice was not shuffled, got %v", shuffled)
	}

	m := make(map[int]bool)
	for _, v := range original {
		m[v] = true
	}
	for _, v := range shuffled {
		if !m[v] {
			t.Errorf("Shuffled slice contains unexpected value %d", v)
		}
	}
}

func TestShuffle_StringSlice(t *testing.T) {
	original := []string{"a", "b", "c", "d", "e"}
	shuffled := make([]string, len(original))
	copy(shuffled, original)

	Shuffle(shuffled)

	if len(original) != len(shuffled) {
		t.Errorf("Expected the shuffled slice length to be %d, got %d", len(original), len(shuffled))
	}

	if reflect.DeepEqual(original, shuffled) {
		t.Errorf("The slice was not shuffled, got %v", shuffled)
	}

	m := make(map[string]bool)
	for _, v := range original {
		m[v] = true
	}
	for _, v := range shuffled {
		if !m[v] {
			t.Errorf("Shuffled slice contains unexpected value %s", v)
		}
	}
}
