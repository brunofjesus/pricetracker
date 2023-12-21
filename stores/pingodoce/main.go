package main

import (
	"fmt"
	"github.com/brunofjesus/pricetracker/stores/connector/dto"
	"log/slog"
	"os"
	"time"

	"github.com/brunofjesus/pricetracker/stores/connector/mq"
	"github.com/brunofjesus/pricetracker/stores/pingodoce/config"
	"github.com/brunofjesus/pricetracker/stores/pingodoce/pkg/crawler"
)

const (
	StoreSlug    = "pingo-doce"
	StoreName    = "Pingo Doce"
	StoreWebSite = "https://mercadao.pt/store/pingo-doce"
)

func main() {
	applicationConfig := config.GetApplicationConfiguration()

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler.WithAttrs([]slog.Attr{
		slog.String("service", "main"),
		slog.String("store", StoreSlug),
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
		Slug:    StoreSlug,
		Name:    StoreName,
		Website: StoreWebSite,
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
		store,
		func(storeProduct dto.StoreProduct) {
			err := publisher.PublishProduct(storeProduct)
			if err != nil {
				fmt.Printf("error: %v", err)
			}
		})
}
