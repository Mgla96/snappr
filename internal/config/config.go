package config

import (
	"fmt"

	"github.com/Mgla96/snappr/internal/adapters/clients"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

// ServicePrefix - imagery producer specific env vars have this prefix
const ServicePrefix = "PR"

// Log level config
type Log struct {
	Level zerolog.Level `default:"info"`
}

type Github struct {
	Token string `required:"true"`
	Owner string `required:"true"`
	Repo  string `required:"true"`
}

type LLM struct {
	Token        string            `required:"true"`
	DefaultModel clients.ModelType `default:"gpt-4-turbo"`
	Endpoint     string
	APIType      clients.APIType `default:"openai"`
	Retries      int             `default:"3"`
}

// Config contains all config parameters for the service
type Config struct {
	Log    Log         `required:"true"`
	Github Github      `required:"true"`
	LLM    LLM         `required:"true"`
	Input  InputConfig `required:"true"`
}

// New returns the parsed config from the environment.
func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(ServicePrefix, cfg); err != nil {
		return nil, fmt.Errorf("error processing envconfig: %w", err)
	}

	return cfg, nil
}
