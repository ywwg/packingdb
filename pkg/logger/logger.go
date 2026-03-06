// Package logger provides leveled logging for packingdb using slog
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// ParseLevel converts a string log level to slog.Level
func ParseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Init initializes the default logger with the specified level
func Init(level string) {
	opts := &slog.HandlerOptions{
		Level: ParseLevel(level),
	}
	handler := slog.NewTextHandler(os.Stderr, opts)
	slog.SetDefault(slog.New(handler))
}

// Convenience functions that use the default logger
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// Fatal logs at error level and exits
func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
