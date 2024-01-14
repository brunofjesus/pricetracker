package product

import (
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/util/list"

	price_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/product"
	stats_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/stats"
	store_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/store"
)

type Updater struct {
	DB                    *sql.DB
	StoreRepository       *store_repository.Repository
	ProductRepository     *product_repository.Repository
	ProductMetaRepository *product_repository.MetaRepository
	PriceRepository       *price_repository.Repository
	StatsRepository       *stats_repository.Repository
}

func (s *Updater) Update(productId int64, storeProduct MqStoreProduct) error {
	tx, err := s.DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Update the product properties
	err = s.ProductRepository.UpdateProduct(
		productId,
		storeProduct.Name,
		storeProduct.Brand,
		storeProduct.ImageLink,
		storeProduct.Link,
		storeProduct.Price,
		storeProduct.Available,
		tx,
	)

	if err != nil {
		return err
	}

	if err = s.updateSkus(productId, storeProduct, tx); err != nil {
		return err
	}

	if err = s.updateEans(productId, storeProduct, tx); err != nil {
		return err
	}

	latestPrice, err := s.PriceRepository.FindLatestPrice(productId, tx)
	if err != nil {
		return err
	}

	if latestPrice.Price != storeProduct.Price ||
		time.Since(latestPrice.DateTime) > time.Hour*4 {

		// Insert price update
		err = s.PriceRepository.CreatePrice(
			productId,
			storeProduct.Price,
			time.Now(),
			tx,
		)

		if err != nil {
			return err
		}

		// Re-calculate the statistics
		err = s.updateStats(productId, storeProduct.Price, tx)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Updater) updateStats(productId int64, latestPrice int, tx *sql.Tx) error {
	prices, err := s.PriceRepository.FindPricesBetween(
		productId,
		time.Now().AddDate(0, 0, -30),
		time.Now(),
		tx,
	)
	if err != nil {
		return err
	}

	if len(prices) == 0 {
		return s.StatsRepository.CreateProductStats(
			productId,
			latestPrice,
			latestPrice,
			latestPrice,
			1,
			decimal.NewFromInt(int64(0)),
			decimal.NewFromInt(int64(0)),
			decimal.NewFromInt(int64(latestPrice)),
			tx,
		)
	}

	minimum := latestPrice
	maximum := latestPrice
	count := len(prices)

	sum := 0

	for _, price := range prices {
		sum += price.Price
		if price.Price < minimum {
			minimum = price.Price
		}
		if price.Price > maximum {
			maximum = price.Price
		}
	}

	average := decimal.NewFromInt(int64(sum)).Div(decimal.NewFromInt(int64(count)))
	difference := decimal.NewFromInt(int64(latestPrice)).Sub(average)
	discountPercent := average.
		Sub(decimal.NewFromInt(int64(latestPrice))).
		Div(average)

	exists, err := s.StatsRepository.HasProductStats(productId, tx)
	if err != nil {
		return err
	}
	if exists {
		err = s.StatsRepository.UpdateProductStats(
			productId,
			latestPrice,
			minimum,
			maximum,
			count,
			difference,
			discountPercent,
			average,
			tx,
		)
	} else {
		err = s.StatsRepository.CreateProductStats(
			productId,
			latestPrice,
			minimum,
			maximum,
			count,
			difference,
			discountPercent,
			average,
			tx,
		)
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *Updater) updateSkus(productId int64, storeProduct MqStoreProduct, tx *sql.Tx) error {
	dbProductSku, err := s.ProductMetaRepository.GetProductSKUs(productId, tx)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// unexpected error
		return err
	}

	if err == nil {
		// update
		previousSkus := list.Map[product_repository.ProductSku, string](dbProductSku, func(m product_repository.ProductSku) string {
			return m.Sku
		})

		toCreate := list.FindMissing[string](storeProduct.SKU, previousSkus)
		toDelete := list.FindMissing[string](previousSkus, storeProduct.SKU)

		if err = s.ProductMetaRepository.CreateSKUs(productId, toCreate, tx); err != nil {
			return err
		}

		if err = s.ProductMetaRepository.DeleteSKUs(productId, toDelete, tx); err != nil {
			return err
		}

	} else if err = s.ProductMetaRepository.CreateSKUs(productId, storeProduct.SKU, tx); err != nil {
		// no records, create
		return err
	}

	return nil
}

func (s *Updater) updateEans(productId int64, storeProduct MqStoreProduct, tx *sql.Tx) error {
	dbProductEan, err := s.ProductMetaRepository.GetProductEANs(productId, tx)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// unexpected error
		return err
	}

	currentEans := filterNumbers(storeProduct.EAN)

	if err == nil {
		// update
		previousEans := list.Map[product_repository.ProductEan, int64](dbProductEan, func(m product_repository.ProductEan) int64 {
			return m.Ean
		})

		toCreate := list.FindMissing[int64](currentEans, previousEans)
		toDelete := list.FindMissing[int64](previousEans, currentEans)

		if err = s.ProductMetaRepository.CreateEANs(productId, toCreate, tx); err != nil {
			return err
		}

		if err = s.ProductMetaRepository.DeleteEANs(productId, toDelete, tx); err != nil {
			return err
		}

	} else if err = s.ProductMetaRepository.CreateEANs(productId, currentEans, tx); err != nil {
		// no records, create
		return err
	}

	return nil
}
