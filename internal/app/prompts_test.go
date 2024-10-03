package app

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewDefaultPromptAndKnowledgeConfig(t *testing.T) {
	// create temp file
	f, err := os.CreateTemp("", "test_snappr_config.yaml")
	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}
	defer os.Remove(f.Name())

	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "NewDefaultPromptAndKnowledgeConfig",
			args: args{
				configPath: f.Name(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewDefaultPromptAndKnowledgeConfig(tt.args.configPath); (err != nil) != tt.wantErr {
				t.Errorf("NewDefaultPromptAndKnowledgeConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			// open file and umarshal to SnapprUserConfig
			f, err := os.ReadFile(tt.args.configPath)
			if err != nil {
				t.Errorf("error reading config file: %v", err)
			}
			err = yaml.Unmarshal(f, &SnapprUserConfig{})
			if err != nil {
				t.Errorf("error unmarshaling SnapprUserConfig: %v", err)
			}

		})
	}
}
