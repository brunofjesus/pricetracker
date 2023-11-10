package mq

import (
	"context"
	"log/slog"
	"time"
)

func SpawnWorker(ctx context.Context, id int) {
	logger := slog.New(
		ctx.Value("logger").(*slog.Logger).
			Handler().WithAttrs(
			[]slog.Attr{
				slog.Int("worker_id", id),
			}),
	)

	logger.Debug("spawning worker")
	for {
		err := newConsumer().Listen(context.WithValue(ctx, "logger", logger))
		if err != nil {
			logger.Error("worker crashed", slog.Any("error", err))
		}

		sleepDuration := time.Millisecond * 5000
		logger.Debug("waiting before re-spawn", slog.Duration("duration", sleepDuration))
		time.Sleep(sleepDuration)
		logger.Debug("re-spawning worker")
	}
}
