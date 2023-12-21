package product

import (
	"database/sql"
	"errors"
	"log/slog"
)

type UrlMatcher struct {
	Finder *Finder
}

func (m *UrlMatcher) Match(storeProduct MqStoreProduct) int64 {
	product, err := m.Finder.FindProductByUrl(storeProduct.Link)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Warn(
			"cannot match URL, this can lead to duplicates",
			slog.Group(
				"match",
				slog.String("store_slug", storeProduct.StoreSlug),
				slog.Any("url", storeProduct.Link),
			),
			slog.Any("error", err),
		)
	} else if err == nil {
		return product.ProductId
	}

	return 0
}

type EanMatcher struct {
	Finder *Finder
}

func (m *EanMatcher) Match(storeProduct MqStoreProduct) int64 {
	id, err := m.Finder.FindProductIdByStoreSlugAndEANs(
		storeProduct.StoreSlug, storeProduct.EAN,
	)

	if err != nil {
		slog.Warn(
			"cannot match EAN, this can lead to duplicates",
			slog.Group(
				"match",
				slog.String("store_slug", storeProduct.StoreSlug),
				slog.Any("eans", storeProduct.EAN),
			),
			slog.Any("error", err),
		)
	}

	return id
}

type SkuMatcher struct {
	Finder *Finder
}

func (m *SkuMatcher) Match(storeProduct MqStoreProduct) int64 {
	id, err := m.Finder.FindProductIdByStoreSlugAndSKUs(
		storeProduct.StoreSlug, storeProduct.SKU,
	)

	if err != nil {
		slog.Warn(
			"cannot match SKU, this can lead to duplicates",
			slog.Group(
				"match",
				slog.String("store_slug", storeProduct.StoreSlug),
				slog.Any("skus", storeProduct.SKU),
			),
			slog.Any("error", err),
		)
	}

	return id
}
