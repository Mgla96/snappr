package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Mgla96/snappr/internal/app"
	"github.com/Mgla96/snappr/internal/config"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var branch, fileRegexPattern, workflowName, repository, repositoryOwner, commitSHA, llmEndpoint string
var printOnly bool
var llmRetries int

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new pull request",
	Long:  `Create a new pull request`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := app.CheckGHToken()
		if err != nil {
			return fmt.Errorf("error no github token: %w", err)
		}
		err = app.CheckLLMToken()
		if err != nil {
			return fmt.Errorf("error no llm token: %w", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
		application := app.SetupNoEnv(&config.Config{
			Log: config.Log{
				Level: zerolog.InfoLevel,
			},
			Github: config.Github{
				Token: os.Getenv("GH_TOKEN"),
				Owner: repositoryOwner,
				Repo:  repository,
			},
			LLM: config.LLM{
				Token:        os.Getenv("LLM_TOKEN"),
				DefaultModel: "gpt-4-turbo",
				Endpoint:     llmEndpoint,
			},
		})
		err := application.ExecuteCreatePR(context.TODO(), commitSHA, branch, workflowName, fileRegexPattern, printOnly)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error executing Create PR")
		}
	},
}

func init() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&repository, "repository", "", "Github repository to discuss, e.g., 'example/repo' (required)")
	createCmd.Flags().StringVar(&repositoryOwner, "repositoryOwner", "", "The account owner of the repository. The name is not case-sensitive. (required)")
	createCmd.Flags().StringVar(&commitSHA, "commitSHA", "", "Commit SHA for creating PR (required)")
	createCmd.Flags().StringVar(&branch, "branch", "", "Branch name to create PR from (required)")
	createCmd.Flags().StringVar(&workflowName, "workflowName", "", "Prompt workflow to use (required)")
	createCmd.Flags().BoolVarP(&printOnly, "printOnly", "p", false, "Only print the resulting PR")
	createCmd.Flags().StringVar(&fileRegexPattern, "fileRegexPattern", `.*\.go$`, "Regex pattern to filter files relevant for PR")
	createCmd.Flags().StringVar(&llmEndpoint, "llmEndpoint", "", "LLM service endpoint (default: 'openai')")
	createCmd.Flags().IntVarP(&llmRetries, "llmRetries", "r", 3, "Number of retries for LLM service calls upon failure")

	setRequiredFlags(createCmd, []string{"repository", "repositoryOwner", "commitSHA", "branch", "workflowName"})
}

func setRequiredFlags(cmd *cobra.Command, flags []string) {
	for _, flag := range flags {
		err := cmd.MarkFlagRequired(flag)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Failed to mark '%s' as required", flag)
		}
	}
}