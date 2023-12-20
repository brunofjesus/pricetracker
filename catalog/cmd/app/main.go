package main

import (
	"context"
	"github.com/brunofjesus/pricetracker/catalog/internal/app"
	"log/slog"
	"os"

	"github.com/brunofjesus/pricetracker/catalog/internal/mq"
	"github.com/brunofjesus/pricetracker/catalog/pkg/rest"
)

func main() {
	logger := newLogger()
	slog.SetDefault(logger)

	logger.Info("starting catalog application")

	appConfig := app.GetApplicationConfiguration()
	environment := loadEnvironment(appConfig)

	logger.Info("will start workers", "workers", appConfig.MessageQueue.ThreadCount)

	workerCtx := context.WithValue(context.Background(), "logger", logger)
	for i := 0; i < appConfig.MessageQueue.ThreadCount; i++ {
		go mq.SpawnWorker(
			workerCtx,
			appConfig,
			environment.Product.Handler,
			environment.Store.Handler,
			i+1,
		)
	}

	// select {}
	rest.ListenAndServe(
		rest.PropsV1{
			ProductFinder: environment.Product.Finder,
			PriceFinder:   environment.Price.Finder,
		},
		8080,
	)
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}).WithAttrs([]slog.Attr{
		slog.String("application", "catalog"),
	}))
}
