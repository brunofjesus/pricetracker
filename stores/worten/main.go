package main

import (
	"log/slog"
	"os"

	"github.com/brunofjesus/pricetracker/stores/worten/definition/catalog"
	"github.com/brunofjesus/pricetracker/stores/worten/integration/mq"
	wortenclient "github.com/brunofjesus/pricetracker/stores/worten/integration/store"
)

const (
	StoreSlug    = "worten"
	StoreName    = "Worten"
	StoreWebSite = "https://www.worten.pt"
)

func main() {

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	publisher, err := mq.NewPublisher(
		slog.New(
			handler.WithAttrs([]slog.Attr{
				slog.String("service", "publisher"),
			}),
		),
	)

	logger := slog.New(handler)
	if err != nil {
		logger.Error("error connecting to MQ", slog.Any("error", err))
		panic(err)
	}

	defer publisher.Close()

	// Register store
	store := catalog.Store{
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
			handler.WithAttrs([]slog.Attr{
				slog.String("service", "productHandler"),
			}),
		),
		Publisher: publisher,
	}

	// Start crawling
	wortenclient.Crawl(
		slog.New(
			handler.WithAttrs([]slog.Attr{
				slog.String("service", "crawler"),
			}),
		),
		productHandler,
	)
}
