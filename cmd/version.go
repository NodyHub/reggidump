/*
Copyright © 2025 Jan Harrie <jan@nody.cc>
*/
package cmd

import (
	"runtime/debug"
	"time"

	"github.com/NodyHub/reggidump/log"
	"github.com/spf13/cobra"
)

var (
	// These variables are populated by the build process
	buildVersion = "dev"
	buildCommit  = "unknown"
	buildDate    = ""
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `The version command prints the version information of the reggidump tool. Includes details such as the version number, commit hash, and build date.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.GetLogger()
		logger.Info("Version information",
			"version", getVersion(),
			"commit", getCommit(),
			"buildDate", getDate())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// getVersion returns the version from build info or the default value
func getVersion() string {
	if buildVersion != "dev" {
		return buildVersion
	}

	// Try to get version from build info
	if info, available := debug.ReadBuildInfo(); available {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return buildVersion
}

// getCommit returns the commit hash from build info or the default value
func getCommit() string {
	if buildCommit != "unknown" {
		return buildCommit
	}

	// Try to get commit from build info
	if info, available := debug.ReadBuildInfo(); available {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return "unknown"
}

// getDate returns the build date or current date if not set
func getDate() string {
	if buildDate != "" {
		return buildDate
	}

	// Try to get time from build info
	if info, available := debug.ReadBuildInfo(); available {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.time" {
				return setting.Value
			}
		}
	}

	// Otherwise use current time
	return time.Now().Format(time.RFC3339)
}
