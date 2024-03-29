package main

import (
	"context"
	"github.com/brunofjesus/pricetracker/catalog/internal/app"
	httpserver "github.com/brunofjesus/pricetracker/catalog/pkg/http"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/frontend"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/rest"
	"github.com/brunofjesus/pricetracker/catalog/pkg/migration"
	"log/slog"
	"os"

	"github.com/brunofjesus/pricetracker/catalog/internal/mq"
)

func main() {
	logger := app.GetLogger()
	slog.SetDefault(logger)

	logger.Info("starting catalog application")

	appConfig := app.GetApplicationConfiguration()
	environment := loadEnvironment(appConfig)

	if appConfig.Database.Migrate {
		logger.Info("running migrations")
		if err := migration.Migrate(environment.DB); err != nil {
			logger.Error("cannot run migration, aborting application startup", slog.Any("error", err))
			os.Exit(2)
		}
	}

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
	host := ""
	if os.Getenv("IS_DEV") == "true" {
		host = "localhost"
		logger.Info("Running in development mode", slog.Any("is_dev", os.Getenv("IS_DEV")))
	}

	httpServerProps := httpserver.ServerProps{
		ApiProps: &rest.V1ApiProps{
			ProductFinder: environment.Product.Finder,
			PriceFinder:   environment.Price.Finder,
		},
		FrontendProps: &frontend.V1FrontendProps{
			ProductFinder: environment.Product.Finder,
			StoreFinder:   environment.Store.Finder,
			PriceFinder:   environment.Price.Finder,
		},
		Host: host,
		Port: 8080,
	}
	err := httpserver.ListenAndServe(httpServerProps)
	if err != nil {
		panic(err)
	}
}
