package main

import (
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/speech/fireworks-admin/internal/app"
	"github.com/speech/fireworks-admin/internal/pkg/logger"
)

func main() {
	application, cleanup, err := app.InitializeApp()
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}
	defer cleanup()

	srv := app.NewServer(application, cleanup)

	if err := srv.Run(); err != nil {
		logger.Error("server run error", slog.Any("error", err))
		os.Exit(1)
	}
}
