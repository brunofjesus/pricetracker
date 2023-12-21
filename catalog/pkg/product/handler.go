package product

import (
	"context"
	"log/slog"
)

type Handler struct {
	Matchers []Matcher
	Creator  *Creator
	Updater  *Updater
}

type Matcher interface {
	Match(storeProduct MqStoreProduct) int64
}

func (s *Handler) Handle(ctx context.Context, storeProduct MqStoreProduct) error {
	logger := ctx.Value("logger").(*slog.Logger)

	var productId int64
	for _, matcher := range s.Matchers {
		if productId = matcher.Match(storeProduct); productId > 0 {
			break
		}
	}

	var err error
	if productId > 0 {
		err = s.Updater.Update(productId, storeProduct)
	} else {
		err = s.Creator.Create(storeProduct)
	}

	if err != nil {
		logger.Error("error on receiver handler", slog.Any("error", err), slog.Any("storeProduct", storeProduct))
	}
	return err
}
