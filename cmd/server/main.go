package main

import (
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/speech/fireworks-admin/internal/app"
	"github.com/speech/fireworks-admin/internal/pkg/logger"
	"github.com/speech/fireworks-admin/internal/pkg/server"
)

func main() {
	application, cleanup, err := app.InitializeApp()
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}
	defer cleanup()

	srv := server.NewServer(application, cleanup)
	app.RegisterRoutes(srv.Echo(), application)

	if err := srv.Start(); err != nil {
		logger.Error("server error", slog.Any("error", err))
		os.Exit(1)
	}
}
