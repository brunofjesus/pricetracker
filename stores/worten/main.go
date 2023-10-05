package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/brunofjesus/pricetracker/stores/worten/config"
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
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)

	// Connect to the MQ
	appConfig := config.GetApplicationConfiguration()
	conn, ch, err := mq.Connect(appConfig.MessageQueue.URL)

	if err != nil {
		logger.Error("error connecting to MQ", slog.Any("error", err))
		panic(err)
	}

	defer conn.Close()
	defer ch.Close()

	// Register store
	store := catalog.Store{
		Slug:    StoreSlug,
		Name:    StoreName,
		Website: StoreWebSite,
	}

	err = mq.PublishStore(ch, store)
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
		MQChannel: ch,
	}

	// Start crawling
	wortenclient.Crawl(
		slog.New(
			handler.WithAttrs([]slog.Attr{
				slog.String("service", "crawler"),
			}),
		),
		time.Second,
		productHandler,
	)
}
