package cmd

import (
	"context"
	"os"

	"github.com/Mgla96/snappr/internal/adapters/clients"
	"github.com/Mgla96/snappr/internal/app"
	"github.com/Mgla96/snappr/internal/config"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// reviewCmd represents the review command
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review a pull request",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return app.CheckTokens()
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
		if commitSHA == "" {
			logger.Fatal().Msg("commitSHA is required")
		}
		if prNumber == 0 {
			logger.Fatal().Msg("prNumber is required and should be greater than 0")
		}

		application := app.SetupNoEnv(&config.Config{
			Input: inputConfig,
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
				DefaultModel: clients.ModelType(llmModel),
				Endpoint:     llmEndpoint,
				APIType:      clients.APIType(llmAPI),
				Retries:      llmRetries,
			},
		})
		err := application.ExecutePRReview(context.Background(), commitSHA, prNumber, workflowName, fileRegexPattern, printOnly)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error executing PR review")
		}
	},
}

func init() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	rootCmd.AddCommand(reviewCmd)

	reviewCmd.Flags().StringVar(&repository, "repository", "", "Github repository to review such as snappr (required)")
	reviewCmd.Flags().StringVar(&repositoryOwner, "repositoryOwner", "", "The account owner of the repository. The name is not case sensitive. (required)")
	reviewCmd.Flags().StringVar(&commitSHA, "commitSHA", "", "Commit SHA to review (required)")
	reviewCmd.Flags().StringVar(&llmModel, "llmModel", "gpt-4-turbo", "Model to use for LLM")
	reviewCmd.Flags().IntVar(&prNumber, "prNumber", 0, "Pull Request number to review (required)")
	reviewCmd.Flags().BoolVarP(&printOnly, "printOnly", "p", false, "Print the review only")
	reviewCmd.Flags().StringVar(&fileRegexPattern, "fileRegexPattern", `.*\.go$`, "Define a regex pattern to filter files to use as context for PR review")
	reviewCmd.Flags().StringVar(&llmEndpoint, "llmEndpoint", "", "Endpoint for the LLM service (defaults to OpenAI)")
	reviewCmd.Flags().IntVarP(&llmRetries, "llmRetries", "r", 3, "Number of retries for LLM API calls when failing to get a valid llm response")
	reviewCmd.Flags().StringVar(&llmAPI, "llmAPI", "openai", "Type of LLM API to use (ollama or openai)")

	err := reviewCmd.MarkFlagRequired("repository")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error marking repository as required")
	}

	err = reviewCmd.MarkFlagRequired("repositoryOwner")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error marking repositoryOwner as required")
	}

	err = reviewCmd.MarkFlagRequired("commitSHA")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error marking commitSHA as required")
	}
	err = reviewCmd.MarkFlagRequired("prNumber")
	if err != nil {
		logger.Err(err).Msg("Error marking prNumber as required")
	}
	reviewCmd.Flags().StringVar(&workflowName, "workflowName", "", "Prompt workflow to use (required)")
	err = reviewCmd.MarkFlagRequired("workflowName")
	if err != nil {
		logger.Err(err).Msg("Error marking workflowName as required")
	}

}
