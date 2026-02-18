package logger

import (
	"log/slog"
	"os"
)

var Logger = slog.Default()

func CreateLogger() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))
}
