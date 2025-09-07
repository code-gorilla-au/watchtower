package logging

import (
	"log/slog"
	"os"
)

type LoggerType string

const (
	LoggerText LoggerType = "Text"
)

var baseLogger *slog.Logger

func init() {
	// start with sane default
	baseLogger = New(slog.LevelInfo)
}

// New logger
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
