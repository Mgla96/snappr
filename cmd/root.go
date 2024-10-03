package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Mgla96/snappr/internal/app"
	"github.com/Mgla96/snappr/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var inputConfig config.InputConfig
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snappr",
	Short: "Save time and catch bugs earlier with snappy PR creation and reviews",
	Long:  `Snappr is a tool for snappy PR creation and review to increase developer velocity. `,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	homeDir := os.Getenv("HOME")
	configDir := filepath.Join(homeDir, ".snappr")
	configFile := filepath.Join(configDir, "config.yaml")
	viper.AddConfigPath(configDir)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err2 := os.Mkdir(configDir, 0755)
		if err2 != nil {
			log.Fatalf("Error creating config directory: %s", err2)
		}
	}

	viper.SetConfigFile(filepath.Join(configDir, "config.yaml"))

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		err = app.NewDefaultPromptAndKnowledgeConfig(configFile)
		if err != nil {
			log.Fatalf("Error creating default config file: %s", err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	if err := viper.Unmarshal(&inputConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "~/.snappr/config.yaml", "config file")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
