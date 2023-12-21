package store

import (
	"log/slog"
	"time"

	"github.com/brunofjesus/pricetracker/stores/worten/config"
)

func Crawl(logger *slog.Logger, handler ProductHandler) {
	appConfig := config.GetApplicationConfiguration()

	categories, err := FindCategories(
		slog.New(
			logger.Handler().WithAttrs([]slog.Attr{
				slog.String("service", "category"),
			}),
		),
	)

	if err != nil {
		logger.Error("Error finding categories", slog.Any("error", err))
		panic(err)
	}

	currentCategoryIdx := 0
	totalCategories := len(categories)

	for categoryId, slug := range categories {
		hasNextPage := true
		page := 0
		currentCategoryIdx = currentCategoryIdx + 1

		logger.Info("Switching category",
			slog.Group("progress",
				slog.Int("current", currentCategoryIdx),
				slog.Int("total", totalCategories),
			),
			slog.Group("category",
				slog.String("id", categoryId),
				slog.String("slug", slug),
			),
		)

		for hasNextPage {
			page = page + 1

			logger.Debug("Getting new page",
				slog.Int("page", page),
				slog.Group("category",
					slog.String("id", categoryId),
					slog.String("slug", slug),
				),
			)

			hasNextPage, err = FindProducts(page, categoryId, slug, handler.Handle)
			if err != nil {
				slog.Error("Error fetching products",
					slog.Int("page", page),
					slog.Group("category",
						slog.String("id", categoryId),
						slog.String("slug", slug),
					),
					slog.Any("error", err),
				)
			}

			time.Sleep(time.Millisecond * time.Duration(appConfig.PolitenessDelayMs))
		}

		logger.Info("Category done",
			slog.Group("category",
				slog.String("id", categoryId),
				slog.String("slug", slug),
			),
			slog.Int("pages", page),
		)
	}
}
