package product

import (
	"database/sql"
	"errors"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/util/list"

	price_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/product"
	store_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/store"
)

type Updater struct {
	DB                    *sql.DB
	StoreRepository       *store_repository.Repository
	ProductRepository     *product_repository.Repository
	ProductMetaRepository *product_repository.MetaRepository
	PriceRepository       *price_repository.Repository
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
	}

	return tx.Commit()
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

	currentEans := filterEANs(storeProduct)

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
