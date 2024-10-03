package app

import (
	"reflect"
	"testing"

	"github.com/Mgla96/snappr/internal/config"
)

func TestGetWorkflowByName(t *testing.T) {
	expectedWorkflow := config.PromptWorkflow{
		Name: "workflow1",
		Steps: []config.PromptWorkflowStep{
			{
				Prompt:      "prompt1",
				InputSource: "text",
			},
		},
	}
	type args struct {
		name         string
		workflowList []config.PromptWorkflow
	}
	tests := []struct {
		name string
		args args
		want *config.PromptWorkflow
	}{
		{
			name: "workflow name not found",
			args: args{
				name:         "workflow1",
				workflowList: []config.PromptWorkflow{{Name: "workflow2"}},
			},
			want: nil,
		},
		{
			name: "workflow name found",
			args: args{
				name: "workflow1",
				workflowList: []config.PromptWorkflow{
					{Name: "workflow2"},
					expectedWorkflow,
				},
			},
			want: &expectedWorkflow,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetWorkflowByName(tt.args.name, tt.args.workflowList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWorkflowByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
