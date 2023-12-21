package main

import (
	"github.com/brunofjesus/pricetracker/stores/connector/dto"
	"log/slog"
	"os"
	"time"

	"github.com/brunofjesus/pricetracker/stores/connector/mq"
	"github.com/brunofjesus/pricetracker/stores/worten/config"
	wortenclient "github.com/brunofjesus/pricetracker/stores/worten/pkg/store"
)

const (
	StoreSlug    = "worten"
	StoreName    = "Worten"
	StoreWebSite = "https://www.worten.pt"
)

func main() {
	applicationConfig := config.GetApplicationConfiguration()

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler.WithAttrs([]slog.Attr{
		slog.String("service", "main"),
		slog.String("store", "worten"),
	}))

	for {
		logger.Info("Running store scrapper")
		run(logger, applicationConfig.MessageQueue.URL)
		logger.Info("Scraping done waiting for next loop", slog.Int64("wait_time_minutes", applicationConfig.LoopIntervalMinutes))
		time.Sleep(time.Minute * time.Duration(applicationConfig.LoopIntervalMinutes))
	}
}

func run(logger *slog.Logger, mqUrl string) {
	publisher, err := mq.NewPublisher(
		slog.New(
			logger.Handler().WithAttrs([]slog.Attr{
				slog.String("service", "publisher"),
			}),
		),
		mqUrl,
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

	// Create the product handler
	var productHandler = wortenclient.ProductHandler{
		Logger: slog.New(
			logger.Handler().WithAttrs([]slog.Attr{
				slog.String("service", "productHandler"),
			}),
		),
		Publisher: publisher,
	}

	// Start crawling
	wortenclient.Crawl(
		slog.New(
			logger.Handler().WithAttrs([]slog.Attr{
				slog.String("service", "crawler"),
			}),
		),
		productHandler,
	)
}
