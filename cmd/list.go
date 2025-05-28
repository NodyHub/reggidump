/*
Copyright © 2025 Jan Harrie <jan@nody.cc>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/NodyHub/reggidump/config"
	"github.com/NodyHub/reggidump/log"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
)

var listConfig config.List
var writer io.Writer

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [registry|file|-]",
	Short: "List available images and tags in a Docker registry",
	Long: `The list command connects to a Docker registry and lists all available images along with their tags.
It can be used to explore the contents of a registry, making it easier to find and manage images.

Input can be provided in three ways:
1. Direct arguments: reggidump list registry1.example.com registry2.example.com
2. From a file: reggidump list targets.txt
3. From stdin: cat targets.txt | reggidump list -`,
	Example: `reggidump list registry.example.com
reggidump list targets.txt
cat targets.txt | reggidump list -`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse common config flags
		listConfig.Common = *ParseCommonConfig(cmd)

		// Get the logger
		logger := log.GetLogger()
		logger.Debug("List command configuration", "config", listConfig)

		// Set up output writer based on output flag
		if listConfig.Output == "stdout" || listConfig.Output == "" {
			writer = os.Stdout
		} else {
			file, err := os.Create(listConfig.Output)
			if err != nil {
				logger.Error("Failed to create output file", "path", listConfig.Output, "error", err)
				os.Exit(1)
			}
			defer file.Close()
			writer = file
			logger.Info("Writing output to file", "path", listConfig.Output)
		}

		// Process targets from arguments
		targets, err := processTargets(args, logger)
		if err != nil {
			logger.Error("Failed to process targets", "error", err)
			os.Exit(1)
		}

		// Process each target
		for _, target := range targets {
			if target == "" || target[0] == '#' {
				// Skip empty lines and comments
				continue
			}
			logger.Info("Listing images for registry", "registry", target)

			// List images for the registry
			if err := listImagesForRegistry(target, listConfig, logger); err != nil {
				logger.Error("Failed to list images", "registry", target, "error", err)
				continue
			}
		}
	},
}

// listImagesForRegistry lists all images for a given registry using google/go-containerregistry
func listImagesForRegistry(registry string, cfg config.List, logger *slog.Logger) error {
	// Parse the registry to ensure proper handling of scheme and URL
	var registryName string

	// Remove scheme if present as name.NewRegistry expects just the registry name without scheme
	registryName = registry
	if strings.HasPrefix(registryName, "http://") {
		registryName = strings.TrimPrefix(registryName, "http://")
	} else if strings.HasPrefix(registryName, "https://") {
		registryName = strings.TrimPrefix(registryName, "https://")
	}

	// Configure the registry client
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	// Create HTTP transport with options
	transportOpts := []remote.Option{}

	// Configure custom user agent
	transportOpts = append(transportOpts, remote.WithUserAgent(cfg.UserAgent))

	// Configure retry
	retryCount := 3 // Default
	if cfg.Rertry != "" {
		if _, err := fmt.Sscanf(cfg.Rertry, "%d", &retryCount); err != nil {
			logger.Warn("invalid retry value, using default", "value", cfg.Rertry, "default", retryCount)
		}
	}

	// Create registry client
	logger.Debug("creating registry client", "registry", registryName)
	reg, err := name.NewRegistry(registryName)
	if err != nil {
		return fmt.Errorf("invalid registry URL: %w", err)
	}

	// Fetch catalog
	logger.Debug("fetching catalog from registry", "registry", reg.Name())
	catalogOpts := append(transportOpts, remote.WithContext(ctx))

	catalog, err := remote.Catalog(ctx, reg, catalogOpts...)
	if err != nil {
		return fmt.Errorf("failed to fetch catalog: %w", err)
	}

	// Filter images if filter is provided
	filteredCatalog := catalog
	if cfg.Filter != "" {
		filteredCatalog = filterImages(catalog, cfg.Filter)
		logger.Debug("applied filter to catalog",
			"filter", cfg.Filter,
			"total", len(catalog),
			"filtered", len(filteredCatalog))
	}

	// Print catalog
	logger.Info("found images in registry", "registry", reg.Name(), "count", len(filteredCatalog))

	// Process each repository
	for _, repo := range filteredCatalog {
		// Get repository reference
		repoRef, err := name.NewRepository(fmt.Sprintf("%s/%s", reg.Name(), repo))
		if err != nil {
			logger.Error("failed to parse repository", "repository", repo, "error", err)
			continue
		}

		// Get tags
		tags, err := remote.List(repoRef, catalogOpts...)
		if err != nil {
			logger.Error("failed to list tags", "repository", repoRef.Name(), "error", err)
			continue
		}

		logger.Info("repository", "name", repoRef.Name(), "tags", len(tags))
		// print server repository and tags in the specified format
		for _, tag := range tags {
			fmt.Fprintf(writer, "%s:%s\n", repoRef.Name(), tag)
		}
	}

	return nil
}

// filterImages applies a filter to the list of images
func filterImages(images []string, filter string) []string {
	if filter == "" {
		return images
	}

	var filtered []string
	for _, img := range images {
		if strings.Contains(img, filter) {
			filtered = append(filtered, img)
		}
	}
	return filtered
}

// processTargets handles the different input methods and returns a slice of targets
func processTargets(args []string, logger *slog.Logger) ([]string, error) {
	var targets []string

	for _, arg := range args {
		switch {
		case arg == "-":
			// Read from stdin
			logger.Debug("Reading targets from stdin")
			stdinTargets, err := readFromReader(os.Stdin, logger)
			if err != nil {
				return nil, fmt.Errorf("failed to read from stdin: %w", err)
			}
			targets = append(targets, stdinTargets...)

		case fileExists(arg):
			// Read from file
			logger.Debug("Reading targets from file", "file", arg)
			file, err := os.Open(arg)
			if err != nil {
				return nil, fmt.Errorf("failed to open file %s: %w", arg, err)
			}
			fileTargets, err := readFromReader(file, logger)
			file.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read from file %s: %w", arg, err)
			}
			targets = append(targets, fileTargets...)

		default:
			// Treat as direct target
			logger.Debug("Adding direct target", "target", arg)
			targets = append(targets, arg)
		}
	}

	return targets, nil
}

// readFromReader reads lines from a reader and returns them as a slice of strings
func readFromReader(r io.Reader, logger *slog.Logger) ([]string, error) {
	var targets []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			targets = append(targets, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return targets, nil
}

// isCommentLine checks if a line starts with a comment character
func isCommentLine(line string) bool {
	return strings.HasPrefix(line, "#")
}

// fileExists checks if a file exists and is not a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add list-specific flags
	listCmd.Flags().StringVarP(&listConfig.Output, "output", "o", "stdout", "Output file path (or 'stdout' for terminal output)")
}
