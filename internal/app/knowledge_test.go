package app

import (
	"os"
	"testing"

	"github.com/Mgla96/snappr/internal/config"
)

func TestRetrieveKnowledge(t *testing.T) {
	file, err := os.CreateTemp("", "test-retrieve-knowledge")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	err = os.WriteFile(file.Name(), []byte("baz"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		sourceName       string
		knowledgeSources []config.KnowledgeSource
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "sourceName not in knowledgeSources list",
			args: args{
				sourceName: "foo",
				knowledgeSources: []config.KnowledgeSource{
					{
						Name:  "bar",
						Type:  config.KnowledgeSourceTypeText,
						Value: "baz",
					},
				},
			},
			want: "",
		},
		{
			name: "sourceName in knowledgeSources list and text",
			args: args{
				sourceName: "foo",
				knowledgeSources: []config.KnowledgeSource{
					{
						Name:  "foo",
						Type:  config.KnowledgeSourceTypeText,
						Value: "baz",
					},
				},
			},
			want: "baz",
		},
		{
			name: "sourceName in knowledgeSources list and file",
			args: args{
				sourceName: "foo",
				knowledgeSources: []config.KnowledgeSource{
					{
						Name:  "foo",
						Type:  config.KnowledgeSourceTypeFile,
						Value: file.Name(),
					},
				},
			},
			want: "baz",
		},
		{
			name: "sourceName in knowledgeSources list and file that doesn't exist",
			args: args{
				sourceName: "foo",
				knowledgeSources: []config.KnowledgeSource{
					{
						Name:  "foo",
						Type:  config.KnowledgeSourceTypeFile,
						Value: "file-does-not-exist",
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "sourceName in knowledgeSources list and URL",
			args: args{
				sourceName: "foo",
				knowledgeSources: []config.KnowledgeSource{
					{
						Name:  "foo",
						Type:  config.KnowledgeSourceTypeURL,
						Value: "baz",
					},
				},
			},
			want: NotImplementedMessage,
		},
		{
			name: "sourceName in knowledgeSources list and API",
			args: args{
				sourceName: "foo",
				knowledgeSources: []config.KnowledgeSource{
					{
						Name:  "foo",
						Type:  config.KnowledgeSourceTypeAPI,
						Value: "baz",
					},
				},
			},
			want: NotImplementedMessage,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RetrieveKnowledge(tt.args.sourceName, tt.args.knowledgeSources)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveKnowledge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RetrieveKnowledge() = %v, want %v", got, tt.want)
			}
		})
	}
}
