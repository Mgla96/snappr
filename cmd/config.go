package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configurations, workflows, and knowledge sources",
	Long: `Manage configurations, workflows, and knowledge sources. 
These configuration values are by default stored in ~/.snappr/config.yaml and can be manually edited.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(listCmd)
}
