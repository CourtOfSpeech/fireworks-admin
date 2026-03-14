package main

import (
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/speech/fireworks-admin/internal/infrastructure/http"
	"github.com/speech/fireworks-admin/pkg/logger"
)

func main() {
	server, err := http.NewServer()
	if err != nil {
		logger.Error("failed to create server", slog.Any("error", err))
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		logger.Error("server error", slog.Any("error", err))
		os.Exit(1)
	}
}
