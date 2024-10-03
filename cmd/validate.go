package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate all configured workflows and knowledge sources",
	Run: func(cmd *cobra.Command, args []string) {
		validateWorkflows()
		validateKnowledgeSources()
	},
}

func validateWorkflows() {
	_ = inputConfig.PromptWorkflows
	fmt.Println("Validation not yet implemented")

}

func validateKnowledgeSources() {
	_ = inputConfig.KnowledgeSources
	fmt.Println("Validation not yet implemented")
}
