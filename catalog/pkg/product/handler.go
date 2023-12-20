package product

import (
	"context"
	"log/slog"
)

type ProductHandler struct {
	ProductMatcher ProductMatcher
	ProductCreator ProductCreator
	ProductUpdater ProductUpdater
}

func (s *ProductHandler) Handle(ctx context.Context, storeProduct MqStoreProduct) error {
	logger := ctx.Value("logger").(*slog.Logger)

	var err error
	if productId := s.ProductMatcher.Match(storeProduct); productId > 0 {
		err = s.ProductUpdater.Update(productId, storeProduct)
	} else {
		err = s.ProductCreator.Create(storeProduct)
	}

	if err != nil {
		logger.Error("error on receiver handler", slog.Any("error", err), slog.Any("storeProduct", storeProduct))
	}
	return err
}
