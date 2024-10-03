package app

import (
	"github.com/Mgla96/snappr/internal/config"
)

var workflows []config.PromptWorkflow

// GetWorkflowByName returns workflow information by name from a list of workflows.
func GetWorkflowByName(name string, workflowList []config.PromptWorkflow) *config.PromptWorkflow {
	for _, wf := range workflowList {
		if wf.Name == name {
			return &wf
		}
	}
	return nil
}
