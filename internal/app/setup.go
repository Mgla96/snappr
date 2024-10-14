package app

import (
	"fmt"
	"os"
	"time"

	"github.com/Mgla96/snappr/internal/adapters/clients"
	"github.com/Mgla96/snappr/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup sets up the application utilizing environment variables.
func Setup() (*App, error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("error setting up config: %w", err)
	}

	zerolog.SetGlobalLevel(cfg.Log.Level)
	ghClient := clients.NewGithubClient(cfg.Github.Token, log.Logger)
	llmClient := clients.NewOpenAIClient(cfg.LLM.Token)

	return New(cfg, ghClient, llmClient, log.Logger), nil
}

// SetupNoEnv sets up the application from a config struct instead of utilizing environment variables.
func SetupNoEnv(cfg *config.Config) *App {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	zerolog.SetGlobalLevel(cfg.Log.Level)

	ghClient := clients.NewGithubClient(cfg.Github.Token, log.Logger)

	var llmClient *clients.OpenAIClient
	if cfg.LLM.Endpoint != "" {
		llmClient = clients.NewCustomOpenAIClient(cfg.LLM.Token, cfg.LLM.Endpoint)
	} else {
		llmClient = clients.NewOpenAIClient(cfg.LLM.Token)
	}

	return New(cfg, ghClient, llmClient, log.Logger)
}
