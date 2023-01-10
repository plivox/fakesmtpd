package cmd

import (
	"fmt"
	"os"

	"fakesmtpd/internal/config"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	globalConfig *config.Config
	flagConfig   string
	flagLevel    string
)

var rootCmd = &cobra.Command{Use: "fakesmtpd"}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().
		StringVarP(&flagConfig, "config", "c", "fakesmtpd.yml", "Config file")
	rootCmd.PersistentFlags().
		StringVar(&flagLevel, "level", "", "Force log level")
}

func initConfig() {
	// Bind an existing flags
	for _, flag := range []string{"config", "level"} {
		viper.BindPFlag(flag, rootCmd.Flags().Lookup(flag))
	}

	// Global config
	globalConfig = config.NewConfig(flagConfig)

	// Use log level flag
	if flagLevel != "" {
		globalConfig.Log.Level = flagLevel
	}

	// Zerolog global logger
	log.Logger = config.NewLogger(globalConfig.Log).
		With().Timestamp().Logger()
}

// Execute execute main command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
