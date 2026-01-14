// Package logging provides structured logging using slog and context-based logger propagation.
package logging

import (
	"log/slog"
	"os"
)

// LoggerType defines the available logger output formats.
type LoggerType string

const (
	// LoggerText represents a text-based logger output format.
	LoggerText LoggerType = "Text"
)

var baseLogger *slog.Logger

func init() {
	// start with sane default
	baseLogger = New(slog.LevelInfo)
}

// New initializes and returns a new text-based logger with the specified level.
func New(level slog.Level) *slog.Logger {
	baseLogger = textHandler(level)
	return baseLogger
}

// textHandler - generate text logger
func textHandler(level slog.Level) *slog.Logger {
	addSource := level == slog.LevelDebug

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
