package mq

import (
	"context"
	"github.com/brunofjesus/pricetracker/catalog/internal/app"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"github.com/brunofjesus/pricetracker/catalog/pkg/store"
	"log/slog"
	"time"
)

func SpawnWorker(
	ctx context.Context,
	appConfig *app.ApplicationConfiguration,
	productHandler *product.Handler,
	storeHandler *store.Handler,
	id int,
) {
	logger := slog.New(
		ctx.Value("logger").(*slog.Logger).
			Handler().WithAttrs(
			[]slog.Attr{
				slog.Int("worker_id", id),
			}),
	)

	logger.Debug("spawning worker")
	for {

		consumer := Consumer{
			productHandler:           productHandler,
			storeHandler:             storeHandler,
			applicationConfiguration: appConfig,
		}

		err := consumer.Listen(context.WithValue(ctx, "logger", logger))
		if err != nil {
			logger.Error("worker crashed", slog.Any("error", err))
		}

		sleepDuration := time.Millisecond * 5000
		logger.Debug("waiting before re-spawn", slog.Duration("duration", sleepDuration))
		time.Sleep(sleepDuration)
		logger.Debug("re-spawning worker")
	}
}
