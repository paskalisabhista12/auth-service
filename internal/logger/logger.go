package logger

import (
	"log/slog"
	"os"
)

func Init() {
	base := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(&trxHandler{
		next: base,
	}).With(
		"pid", os.Getpid(),
	)

	slog.SetDefault(logger)
}
