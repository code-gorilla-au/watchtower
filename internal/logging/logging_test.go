package logging

import (
	"context"
	"log/slog"
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestLogging(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("LoggerType constants are defined correctly", func(t *testing.T) {
			odize.AssertEqual(t, LoggerText, LoggerType("Text"))
		}).
		Test("New creates a logger with correct level", func(t *testing.T) {
			logger := New(slog.LevelInfo)
			odize.AssertFalse(t, logger == nil)

			debugLogger := New(slog.LevelDebug)
			odize.AssertFalse(t, debugLogger == nil)

			errorLogger := New(slog.LevelError)
			odize.AssertFalse(t, errorLogger == nil)
		}).
		Test("textHandler creates logger with correct settings", func(t *testing.T) {
			logger := textHandler(slog.LevelInfo)
			odize.AssertFalse(t, logger == nil)

			debugLogger := textHandler(slog.LevelDebug)
			odize.AssertFalse(t, debugLogger == nil)
		}).
		Test("init initializes baseLogger", func(t *testing.T) {
			odize.AssertFalse(t, baseLogger == nil)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestContext(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("FromContext returns baseLogger when no logger in context", func(t *testing.T) {
			ctx := context.Background()
			logger := FromContext(ctx)
			odize.AssertFalse(t, logger == nil)
			odize.AssertEqual(t, logger, baseLogger)
		}).
		Test("FromContext returns logger from context when present", func(t *testing.T) {
			ctx := context.Background()
			testLogger := New(slog.LevelDebug)
			ctxWithLogger := WithContext(ctx, testLogger)

			retrievedLogger := FromContext(ctxWithLogger)
			odize.AssertFalse(t, retrievedLogger == nil)
			odize.AssertEqual(t, retrievedLogger, testLogger)
		}).
		Test("WithContext attaches logger to context", func(t *testing.T) {
			ctx := context.Background()
			testLogger := New(slog.LevelWarn)

			newCtx := WithContext(ctx, testLogger)
			odize.AssertFalse(t, newCtx == nil)

			retrievedLogger := FromContext(newCtx)
			odize.AssertEqual(t, retrievedLogger, testLogger)
		}).
		Test("FromContext returns baseLogger for invalid context value", func(t *testing.T) {
			ctx := context.WithValue(context.Background(), loggerContextKey{}, "not a logger")
			logger := FromContext(ctx)
			odize.AssertEqual(t, logger, baseLogger)
		}).
		Run()

	odize.AssertNoError(t, err)
}
