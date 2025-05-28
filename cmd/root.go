/*
Copyright © 2025 Jan Harrie <jan@nody.cc>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/NodyHub/reggidump/config"
	"github.com/NodyHub/reggidump/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var commonConfig config.Common
var verbose bool

// Version information
var (
	version = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "reggidump",
	Short: "Inspect and dump Docker images from a registry",
	Long: `reggidump is a tool for inspecting and dumping Docker images from a registry.
It allows users to list available images, download image layers, and inspect image manifests.
It is designed to work with Docker registries, providing a simple command-line interface for managing images.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize logger
		log.InitLogger(commonConfig.LogLevel, commonConfig.LogFile, verbose)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.reggidump.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose output")

	// Common flags
	rootCmd.PersistentFlags().StringVar(&commonConfig.Auth, "auth", "", "authentication credentials in format 'username:password'")
	rootCmd.PersistentFlags().StringVar(&commonConfig.Filter, "filter", "", "filter pattern for images")
	rootCmd.PersistentFlags().StringSliceVar(&commonConfig.Header, "header", []string{}, "additional HTTP headers to send (can be used multiple times)")
	rootCmd.PersistentFlags().StringVar(&commonConfig.LogLevel, "log-level", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVar(&commonConfig.LogFile, "log-file", "stderr", "log file path (or 'stderr'/'stdout')")
	rootCmd.PersistentFlags().StringVar(&commonConfig.Rertry, "retry", "3", "number of retries for HTTP requests")
	rootCmd.PersistentFlags().IntVar(&commonConfig.Timeout, "timeout", 1, "HTTP timeout in seconds")
	rootCmd.PersistentFlags().StringVar(&commonConfig.UserAgent, "user-agent", "reggidump "+version, "user agent string for HTTP requests")
}

// ParseCommonConfig parses the common flags from cobra command
func ParseCommonConfig(cmd *cobra.Command) *config.Common {
	// Create a copy of the common config
	parsedConfig := commonConfig

	// Bind flags to viper config if needed
	// This could be extended if we want to support config files more robustly

	return &parsedConfig
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".reggidump" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".reggidump")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
