package main

import (
	"context"
	"github.com/brunofjesus/pricetracker/catalog/internal/app"
	httpserver "github.com/brunofjesus/pricetracker/catalog/pkg/http"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/frontend"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/rest"
	"log/slog"
	"os"

	"github.com/brunofjesus/pricetracker/catalog/internal/mq"
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
	httpServerProps := httpserver.ServerProps{
		ApiProps: &rest.V1ApiProps{
			ProductFinder: environment.Product.Finder,
			PriceFinder:   environment.Price.Finder,
		},
		FrontendProps: &frontend.V1FrontendProps{},
		Port:          8080,
	}
	err := httpserver.ListenAndServe(httpServerProps)
	if err != nil {
		panic(err)
	}
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}).WithAttrs([]slog.Attr{
		slog.String("application", "catalog"),
	}))
}
