package main

import (
	"log/slog"
	"os"
)

// Log
func main() {
	// h := slog.NewTextHandler(os.Stderr, nil)
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	log := slog.New(h).With("app", "exp")
	// log := slog.New(h).WithGroup("app")

	log.Debug("debug message")
	log.Info("ifo message", "request_id", 1, "user", "Luís")
	log.Warn("warm message")
	log.Error("error message", "request_id", 1)
}
