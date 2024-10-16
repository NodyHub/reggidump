package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/NodyHub/reggidump/registry"
	"github.com/alecthomas/kong"
)

type CLI struct {
	Auth         string `short:"a" name:"auth" help:"Registry authentication header."`
	Dump         string `short:"d" name:"dump" optional:"" help:"Dump image layers to specified directory."`
	FailCount    int    `short:"f" name:"fail-count" default:"5" help:"Number of failed downloads before giving up a server."`
	List         bool   `short:"l" name:"list" default:"false" help:"List available images+tags."`
	ManifestOnly bool   `short:"m" name:"manifest-only" default:"false" help:"Dump only image manifest."`
	Output       string `short:"o" name:"output" optional:"" default:"-" help:"Output file. Default is stdout."`
	Parallel     int    `short:"P" name:"parallel" default:"5" help:"Number of parallel downloads."`
	Ping         bool   `short:"p" name:"ping" default:"false" help:"Check if target is a registry"`
	Retry        int    `short:"r" name:"retry" default:"5" help:"Number of retries a download."`
	Timeout      int    `short:"t" name:"timeout" optional:"" default:"5" help:"Timeout in seconds for registry operations."`
	UserAgent    string `short:"u" name:"user-agent" default:"reggidump" help:"User agent string."`
	Verbose      bool   `short:"v" name:"verbose" help:"Enable verbose output."`
	Version      kong.VersionFlag
	Target       []string `arg:"" name:"target" help:"Dump targets. Can be a file, a registry address or - for stdin."`
}

// version info
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// main function
func main() {

	var cli CLI
	kong.Parse(&cli,
		kong.Description("Dump all Docker images from a registry."),
		kong.UsageOnError(),
		kong.Vars{
			"version": fmt.Sprintf("%s (%s), commit %s, built at %s", filepath.Base(os.Args[0]), version, commit, date),
		},
	)

	// Check for verbose output
	logLevel := slog.LevelError
	if cli.Verbose {
		logLevel = slog.LevelDebug
	}

	// setup logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	// set http timeout
	registry.SetHttpTimeout(cli.Timeout)

	// set retry
	registry.SetHttpRetry(cli.Retry)

	// set auth token
	if len(cli.Auth) > 0 {
		registry.SetAuthToken(cli.Auth)
	}

	// set user agent
	if len(cli.UserAgent) > 0 {
		registry.SetUserAgent(cli.UserAgent)
	}

	// check for output destination
	var output *os.File
	if cli.Output == "-" {
		output = os.Stdout
	} else {
		var err error
		output, err = os.Create(cli.Output)
		if err != nil {
			logger.Error("failed to create output file", "err", err)
			os.Exit(1)
		}
		defer output.Close()
	}

	// check if targets is a file
	targets, err := parseTarget(cli.Target)
	if err != nil {
		logger.Error("failed to parse target", "err", err)
		os.Exit(1)
	}

	// iterate over targets
	for _, t := range targets {
		if len(t) == 0 {
			continue
		}
		if strings.HasPrefix(t, "#") {
			continue
		}
		logger.Debug("processing next", "target", t)
		s := registry.NewServer(t)

		// ping server to check if it is a registry and negotiate protocol + port
		if cli.Ping {
			if err := s.Ping(); err == nil {
				logger.Info("server is a registry", "address", s.Address)
				fmt.Fprintln(output, s.GetService())
			} else {
				logger.Info("server is not a registry", "address", s.Address, "err", err)
				continue
			}
		}

		// print available images
		if cli.List {

			// fetch image list + tags
			logger.Debug("list images and tags", "server", s.Address)
			if err := s.FetchImagesAndTags(logger); err != nil {
				logger.Error("failed to fetch images", "err", err)
				continue
			}

			// print images
			for _, i := range s.Images {
				for _, tag := range i.Tags {
					if _, err := fmt.Fprintf(output, "%s/%s:%s\n", s.GetService(), i.Name, tag.Name); err != nil {
						logger.Error("failed to write to output", "err", err)
						os.Exit(1)
					}
				}
			}
		}

		// dump image layers
		if len(cli.Dump) > 0 {
			logger.Debug("start dump", "path", cli.Dump, "server", s.Address)
			if err := s.Dump(logger, cli.Dump, cli.ManifestOnly, int32(cli.FailCount), cli.Parallel); err != nil {
				logger.Error("failed, that sucks!", "err", err)
			}
		}

	}
}

func parseTarget(input []string) ([]string, error) {

	var targets []string

	for _, i := range input {
		if len(i) == 0 {
			return nil, fmt.Errorf("empty target")
		}

		if i == "-" {
			// read from stdin
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				targets = append(targets, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return nil, err
			}
			return targets, nil
		}

		// check if target is a file
		if _, err := os.Stat(i); err == nil {
			// read all lines in file
			c, err := os.ReadFile(i)
			if err != nil {
				return nil, err
			}
			targets = append(targets, strings.Split(string(c), "\n")...)
			continue
		}

		// add input as target
		targets = append(targets, i)

	}

	return targets, nil
}
