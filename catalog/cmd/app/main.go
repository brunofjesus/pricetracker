package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/brunofjesus/pricetracker/catalog/config"
	"github.com/brunofjesus/pricetracker/catalog/internal/mq"
	"github.com/brunofjesus/pricetracker/catalog/pkg/rest"
)

func main() {
	logger := newLogger()
	slog.SetDefault(logger)

	logger.Info("starting catalog application")

	appConfig := config.GetApplicationConfiguration()

	logger.Info("will start workers", "workers", appConfig.MessageQueue.ThreadCount)

	workerCtx := context.WithValue(context.Background(), "logger", logger)
	for i := 0; i < appConfig.MessageQueue.ThreadCount; i++ {
		go mq.SpawnWorker(workerCtx, i+1)
	}

	// select {}
	rest.ListenAndServe(8080)
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}).WithAttrs([]slog.Attr{
		slog.String("application", "catalog"),
	}))
}
