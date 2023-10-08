package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/brunofjesus/pricetracker/stores/pingodoce/config"
	"github.com/brunofjesus/pricetracker/stores/pingodoce/crawler"
	"github.com/brunofjesus/pricetracker/stores/pingodoce/definition"
	"github.com/brunofjesus/pricetracker/stores/pingodoce/mq"
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
		run(logger)
		logger.Info("Scraping done waiting for next loop", slog.Int64("wait_time_ms", applicationConfig.LoopIntervalMs))
		time.Sleep(time.Millisecond * time.Duration(applicationConfig.LoopIntervalMs))
	}
}

func run(logger *slog.Logger) {
	publisher, err := mq.NewPublisher(
		slog.New(
			logger.Handler().WithAttrs([]slog.Attr{
				slog.String("service", "publisher"),
			}),
		),
	)

	if err != nil {
		logger.Error("error connecting to MQ", slog.Any("error", err))
		panic(err)
	}

	defer publisher.Close()

	// Register store
	store := definition.Store{
		Slug:    StoreSlug,
		Name:    StoreName,
		Website: StoreWebSite,
	}

	err = publisher.PublishStore(store)
	if err != nil {
		logger.Error("error publishing store to MQ", slog.Any("error", err))
		panic(err)
	}

	crawler.Crawl(store, func(storeProduct definition.StoreProduct) {
		err := publisher.PublishProduct(storeProduct)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
	})
}
