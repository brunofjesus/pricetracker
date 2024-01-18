package main

import (
	"fmt"
	"github.com/brunofjesus/pricetracker/stores/connector/dto"
	"github.com/brunofjesus/pricetracker/stores/connector/mq"
	"github.com/brunofjesus/pricetracker/stores/example/config"
	"github.com/brunofjesus/pricetracker/stores/example/pkg/crawler"
	"log/slog"
	"os"
	"time"
)

const (
	StoreSlug    = "example-store"
	StoreName    = "Example Store"
	StoreWebSite = "http://localhost"
)

func main() {
	// Get the configurations stored in the config.yaml file
	applicationConfig := config.GetApplicationConfiguration()

	logger := logger()
	for {
		logger.Info("Running example store price fetcher")
		if err := run(logger, applicationConfig); err != nil {
			logger.Error("scrapper failed", slog.Any("error", err))
		}
		logger.Info("Scraping done waiting for next loop", slog.Int64("wait_time_minutes", applicationConfig.LoopIntervalMinutes))
		time.Sleep(time.Minute * time.Duration(applicationConfig.LoopIntervalMinutes))
	}
}

// run connects to the Message Queue, registers the store in the catalog and then starts
// looking for products on the store.
func run(logger *slog.Logger, appConfig *config.ApplicationConfiguration) error {
	publisher, err := mq.NewPublisher(logger, appConfig.MessageQueue.URL)

	if err != nil {
		return fmt.Errorf("cannot connect to MQ, %w", err)
	}

	defer publisher.Close()

	if err := registerStore(logger, publisher); err != nil {
		return fmt.Errorf("cannot register store through MQ, %w", err)
	}

	// Execute the crawler and provide an anonymous callback function which publishes the product to
	// the catalog and logs errors (if any)
	crawler.Crawl(
		func(storeProduct dto.StoreProduct) {
			if err := publisher.PublishProduct(storeProduct); err != nil {
				logger.Error("error publishing product", slog.Any("error", err))
			}
		},
	)

	return nil
}

// registers the store into the catalog
func registerStore(logger *slog.Logger, publisher *mq.Publisher) error {
	store := dto.Store{
		Slug:    StoreSlug,
		Name:    StoreName,
		Website: StoreWebSite,
	}

	logger.Info("Registering store into the catalog", slog.Any("store", store))

	return publisher.PublishStore(store)
}

// Creates a logger
func logger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler.WithAttrs([]slog.Attr{
		slog.String("store", StoreSlug),
	}))
}
