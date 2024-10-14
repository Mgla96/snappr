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

var branch string
var fileRegexPattern string
var llmRetries int
var llmModel string
var llmAPI string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new pull request",
	Long:  `Create a new pull request`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return app.CheckTokens()
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
				DefaultModel: clients.ModelType(llmModel),
				Endpoint:     llmEndpoint,
				APIType:      clients.APIType(llmAPI),
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

	createCmd.Flags().StringVar(&repository, "repository", "", "Github repository to review such as snappr (required)")
	createCmd.Flags().StringVar(&repositoryOwner, "repositoryOwner", "", "The account owner of the repository. The name is not case sensitive. (required)")
	createCmd.Flags().StringVar(&commitSHA, "commitSHA", "", "Commit SHA to create PR from (required)")
	createCmd.Flags().StringVar(&branch, "branch", "", "Branch name to create PR from (required)")
	createCmd.Flags().StringVar(&llmModel, "llmModel", "gpt-4-turbo", "Model to use for LLM")
	createCmd.Flags().BoolVarP(&printOnly, "printOnly", "p", false, "Print the created PR only")
	createCmd.Flags().StringVar(&fileRegexPattern, "fileRegexPattern", `.*\.go$`, "Define a regex pattern to filter files to use as context for PR creation")
	createCmd.Flags().StringVar(&llmEndpoint, "llmEndpoint", "", "Endpoint for the LLM service (defaults to OpenAI)")
	createCmd.Flags().IntVarP(&llmRetries, "llmRetries", "r", 3, "Number of retries for LLM API calls when failing to get a valid llm response")
	createCmd.Flags().StringVar(&llmAPI, "llmAPI", "openai", "Type of LLM API to use (ollama or openai)")
	createCmd.Flags().StringVar(&workflowName, "workflowName", "", "Prompt workflow to use (required)")

	err := createCmd.MarkFlagRequired("repository")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error marking repository as required")
	}

	err = createCmd.MarkFlagRequired("repositoryOwner")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error marking repositoryOwner as required")
	}

	err = createCmd.MarkFlagRequired("commitSHA")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error marking commitSHA as required")
	}
	err = createCmd.MarkFlagRequired("branch")
	if err != nil {
		logger.Err(err).Msg("Error marking branch as required")
	}

	err = createCmd.MarkFlagRequired("workflowName")
	if err != nil {
		logger.Err(err).Msg("Error marking workflowName as required")
	}

}
