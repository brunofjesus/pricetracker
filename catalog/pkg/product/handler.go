package product

import (
	"context"
	"errors"
	"log/slog"
)

type Handler struct {
	Matchers []Matcher
	Creator  *Creator
	Updater  *Updater
}

type MqStoreProduct struct {
	StoreSlug string
	EAN       []string
	SKU       []string
	Name      string
	Brand     string
	Price     int
	Available bool
	ImageLink string
	Link      string
	Currency  string
}

type Matcher interface {
	Match(storeProduct MqStoreProduct) int64
}

func (s *Handler) Handle(ctx context.Context, storeProduct MqStoreProduct) error {
	logger := ctx.Value("logger").(*slog.Logger)

	var err = validate(storeProduct)
	if err != nil {
		logger.Error("invalid product received", slog.Any("error", err), slog.Any("storeProduct", storeProduct))
		return err
	}

	var productId int64
	for _, matcher := range s.Matchers {
		if productId = matcher.Match(storeProduct); productId > 0 {
			break
		}
	}

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

func validate(product MqStoreProduct) error {
	if product.Price < 0 {
		return errors.New("price can't be negative")
	}

	if len(product.Currency) == 0 {
		return errors.New("currency should be in ISO 4217 format")
	}

	if len(product.Name) == 0 {
		return errors.New("product name is required")
	}

	if len(product.Link) == 0 {
		return errors.New("product url is required")
	}

	if len(product.StoreSlug) == 0 {
		return errors.New("store slug is required")
	}

	return nil
}
