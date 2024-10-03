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

var repository string
var repositoryOwner string
var commitSHA string
var prNumber int
var printOnly bool
var llmEndpoint string
var workflowName string

// reviewCmd represents the review command
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review a pull request",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:
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
		if commitSHA == "" {
			logger.Fatal().Msg("commitSHA is required")
		}
		if prNumber == 0 {
			logger.Fatal().Msg("prNumber is required and should be greater than 0")
		}

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
		err := application.ExecutePRReview(context.Background(), commitSHA, prNumber, workflowName, printOnly)
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
	reviewCmd.Flags().IntVar(&prNumber, "prNumber", 0, "Pull Request number to review (required)")
	reviewCmd.Flags().BoolVarP(&printOnly, "printOnly", "p", false, "Print the review only")
	reviewCmd.Flags().StringVar(&llmEndpoint, "llmEndpoint", "", "Endpoint for the LLM service (defualts to OpenAI)")
	reviewCmd.Flags().IntVarP(&llmRetries, "llmRetries", "r", 3, "Number of retries for LLM API calls when failing to get a valid llm response")

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
