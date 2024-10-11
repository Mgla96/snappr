package clients

import (
	"context"
	"testing"

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
		// TODO: Add test cases.
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
