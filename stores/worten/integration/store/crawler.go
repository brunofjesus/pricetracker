package store

import (
	"log/slog"
	"time"
)

func Crawl(logger *slog.Logger, pageSwitchDelay time.Duration, handler ProductHandler) {

	categories, err := FindCategories()
	if err != nil {
		logger.Error("Error finding categories", slog.Any("error", err))
		panic(err)
	}

	for categoryId, slug := range categories {
		hasNextPage := true
		page := 0

		logger.Debug("Switching category",
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

			time.Sleep(pageSwitchDelay)
		}

		logger.Debug("Switching category",
			slog.Group("category",
				slog.String("id", categoryId),
				slog.String("slug", slug),
			),
			slog.Int("pages", page),
		)
	}
}
