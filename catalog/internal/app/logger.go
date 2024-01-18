package app

import (
	"log/slog"
	"os"
	"sync"
)

var loggerOnce sync.Once
var loggerInstance *slog.Logger

func GetLogger() *slog.Logger {
	loggerOnce.Do(func() {
		loggerInstance = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}).WithAttrs([]slog.Attr{
			slog.String("application", "catalog"),
		}))
	})
	return loggerInstance
}
