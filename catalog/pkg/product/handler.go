package product

import (
	"context"
	"log/slog"
)

type Handler struct {
	Matcher *Matcher
	Creator *Creator
	Updater *Updater
}

func (s *Handler) Handle(ctx context.Context, storeProduct MqStoreProduct) error {
	logger := ctx.Value("logger").(*slog.Logger)

	var err error
	if productId := s.Matcher.Match(storeProduct); productId > 0 {
		err = s.Updater.Update(productId, storeProduct)
	} else {
		err = s.Creator.Create(storeProduct)
	}

	if err != nil {
		logger.Error("error on receiver handler", slog.Any("error", err), slog.Any("storeProduct", storeProduct))
	}
	return err
}
