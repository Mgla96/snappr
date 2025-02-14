package cmd

import (
	"fmt"

	"github.com/Mgla96/snappr/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured workflows and knowledge sources",
	Run: func(cmd *cobra.Command, args []string) {
		listWorkflows(inputConfig.PromptWorkflows)
		listKnowledgeSources(inputConfig.KnowledgeSources)
	},
}

func listWorkflows(workflows []config.PromptWorkflow) {
	fmt.Println("Configured Workflows:")
	for _, wf := range workflows {
		fmt.Printf("- %s\n", wf.Name)
	}
}

func listKnowledgeSources(knowledgeSources []config.KnowledgeSource) {
	fmt.Println("\nConfigured Knowledge Sources:")
	for _, ks := range knowledgeSources {
		fmt.Printf("- %s (%s)\n", ks.Name, ks.Type)
	}
}
