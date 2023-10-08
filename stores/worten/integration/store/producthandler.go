package store

import (
	"log/slog"

	"github.com/brunofjesus/pricetracker/stores/worten/definition/catalog"
	"github.com/brunofjesus/pricetracker/stores/worten/definition/store"
	"github.com/brunofjesus/pricetracker/stores/worten/integration/mapping"
	"github.com/brunofjesus/pricetracker/stores/worten/integration/mq"
)

type ProductHandler struct {
	Logger    *slog.Logger
	Publisher *mq.Publisher
}

func (h *ProductHandler) Handle(wph store.WortenProductHit) error {
	var storeProduct catalog.StoreProduct = catalog.StoreProduct{}
	err := mapping.MapWortenProductToCatalogProduct(wph, &storeProduct)
	if err != nil {
		h.Logger.Error("Error mapping to store product", slog.String("service", "mapper"), slog.Any("error", err))
		return err
	}

	err = h.Publisher.PublishProduct(storeProduct)
	if err != nil {
		h.Logger.Error("Error publishing product", slog.Any("error", err))
	}

	return err
}
