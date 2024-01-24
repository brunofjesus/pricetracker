package main

import (
	"fmt"
	"github.com/brunofjesus/pricetracker/stores/connector/dto"
	"github.com/brunofjesus/pricetracker/stores/connector/mq"
	"github.com/brunofjesus/pricetracker/stores/motocard/config"
	"github.com/brunofjesus/pricetracker/stores/motocard/pkg/crawler"
	"log/slog"
	"os"
	"time"
)

func main() {
	applicationConfig := config.GetApplicationConfiguration()

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler.WithAttrs([]slog.Attr{
		slog.String("service", "main"),
		slog.String("store", applicationConfig.Store.Slug),
	}))

	for {
		logger.Info("Running store scrapper")
		run(logger, applicationConfig)
		logger.Info("Scraping done waiting for next loop", slog.Int64("wait_time_minutes", applicationConfig.LoopIntervalMinutes))
		time.Sleep(time.Minute * time.Duration(applicationConfig.LoopIntervalMinutes))
	}
}

func run(logger *slog.Logger, appConfig *config.ApplicationConfiguration) {
	publisher, err := mq.NewPublisher(
		slog.New(
			logger.Handler().WithAttrs([]slog.Attr{
				slog.String("service", "publisher"),
			}),
		),
		appConfig.MessageQueue.URL,
	)

	if err != nil {
		logger.Error("error connecting to MQ", slog.Any("error", err))
		panic(err)
	}

	defer publisher.Close()

	// Register store
	store := dto.Store{
		Slug:    appConfig.Store.Slug,
		Name:    appConfig.Store.Name,
		Website: fmt.Sprintf("https://www.motocard.com/%s/", appConfig.Store.Country),
	}

	err = publisher.PublishStore(store)
	if err != nil {
		logger.Error("error publishing store to MQ", slog.Any("error", err))
		panic(err)
	}

	crawler.Crawl(
		slog.New(
			logger.Handler().WithAttrs([]slog.Attr{
				slog.String("service", "crawler"),
			}),
		),
		appConfig.PolitenessDelayMs,
		func(storeProduct dto.StoreProduct) {
			err := publisher.PublishProduct(storeProduct)
			if err != nil {
				fmt.Printf("error: %v", err)
			}
		})
}
