package log

import (
	"io"
	"log/slog"
	"os"
)

var logger *slog.Logger

// InitLogger initializes the global logger with the specified level and output
func InitLogger(level string, output string, verbose bool) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// If verbose is set, override level to debug
	if verbose {
		logLevel = slog.LevelDebug
	}

	var writer io.Writer
	if output == "stderr" || output == "" {
		writer = os.Stderr
	} else if output == "stdout" {
		writer = os.Stdout
	} else {
		// Try to open a file
		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// Fall back to stderr on error
			writer = os.Stderr
		} else {
			writer = file
		}
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.NewTextHandler(writer, opts)
	logger = slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

// GetLogger returns the global logger
func GetLogger() *slog.Logger {
	if logger == nil {
		return InitLogger("info", "stderr", false)
	}
	return logger
}
