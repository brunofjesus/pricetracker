package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/config"
	"github.com/brunofjesus/pricetracker/catalog/src/service/mq"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}).WithAttrs([]slog.Attr{
		slog.String("application", "catalog"),
	}))
	slog.SetDefault(logger)

	ctx := context.WithValue(context.Background(), "logger", logger)

	logger.Info("starting catalog application")

	appConfig := config.GetApplicationConfiguration()

	logger.Info("will start workers", "workers", appConfig.MessageQueue.ThreadCount)

	for i := 0; i < appConfig.MessageQueue.ThreadCount; i++ {
		go worker(ctx, i+1)
	}

	select {}
}

func worker(ctx context.Context, id int) {
	logger := slog.New(
		ctx.Value("logger").(*slog.Logger).
			Handler().WithAttrs(
			[]slog.Attr{
				slog.Int("worker_id", id),
			}),
	)

	logger.Debug("spawning worker")
	for {
		err := mq.NewConsumer().Listen(context.WithValue(ctx, "logger", logger))
		if err != nil {
			logger.Error("worker crashed", slog.Any("error", err))
		}

		sleepDuration := time.Millisecond * 5000
		logger.Debug("waiting before re-spawn", slog.Duration("duration", sleepDuration))
		time.Sleep(sleepDuration)
		logger.Debug("re-spawning worker")
	}
}
