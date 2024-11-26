package helpers

import (
	"log/slog"
	"os"
)

func SetupLogger(level string, loggingType string) *slog.Logger {
	var log *slog.Logger

	switch loggingType {
	case "text":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: loggingLevel(level)}))

	case "json":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: loggingLevel(level)}))

	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: loggingLevel(level)}))
	}

	return log
}

func loggingLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
