package updater

import (
	"database/sql"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"
	"github.com/brunofjesus/pricetracker/catalog/src/model"
	"github.com/brunofjesus/pricetracker/catalog/src/repository"
	"github.com/brunofjesus/pricetracker/catalog/src/util/list"

	price_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product"
	product_meta_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product/meta"
	store_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/store"
)

var once sync.Once
var instance ProductUpdater

type ProductUpdater interface {
	Update(productId int64, storeProduct datasource.StoreProduct) error
}

type productUpdater struct {
	db                    *sql.DB
	storeRepository       store_repository.StoreRepository
	productRepository     product_repository.ProductRepository
	productMetaRepository product_meta_repository.ProductMetaRepository
	priceRepository       price_repository.PriceRepository
}

func GetProductUpdater() ProductUpdater {
	once.Do(func() {
		instance = &productUpdater{
			db:                    repository.GetDatabaseConnection(),
			storeRepository:       store_repository.GetStoreRepository(),
			productRepository:     product_repository.GetProductRepository(),
			productMetaRepository: product_meta_repository.GetProductMetaRepository(),
			priceRepository:       price_repository.GetPriceRepository(),
		}
	})
	return instance
}

// Update implements ProductUpdater.
func (s *productUpdater) Update(productId int64, storeProduct datasource.StoreProduct) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Update the product properties
	err = s.productRepository.UpdateProduct(
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

	latestPrice, err := s.priceRepository.GetLatestPrice(productId, tx)
	if err != nil {
		return err
	}

	if !latestPrice.Price.Equal(storeProduct.Price) ||
		time.Since(latestPrice.DateTime) > time.Hour {

		// Insert price update
		err = s.priceRepository.CreatePrice(
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

func (s *productUpdater) updateSkus(productId int64, storeProduct datasource.StoreProduct, tx *sql.Tx) error {
	dbProductSku, err := s.productMetaRepository.GetProductSKUs(productId, tx)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// unexpected error
		return err
	}

	if err == nil {
		// update
		previousSkus := list.Map[model.ProductSku, string](dbProductSku, func(m model.ProductSku) string {
			return m.Sku
		})

		toCreate := list.FindMissing[string](storeProduct.SKU, previousSkus)
		toDelete := list.FindMissing[string](previousSkus, storeProduct.SKU)

		if err = s.productMetaRepository.CreateSKUs(productId, toCreate, tx); err != nil {
			return err
		}

		if err = s.productMetaRepository.DeleteSKUs(productId, toDelete, tx); err != nil {
			return err
		}

	} else if err = s.productMetaRepository.CreateSKUs(productId, storeProduct.SKU, tx); err != nil {
		// no records, create
		return err
	}

	return nil
}

func (s *productUpdater) updateEans(productId int64, storeProduct datasource.StoreProduct, tx *sql.Tx) error {
	dbProductEan, err := s.productMetaRepository.GetProductEANs(productId, tx)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// unexpected error
		return err
	}

	currentEans := filterEANs(storeProduct)

	if err == nil {
		// update
		previousEans := list.Map[model.ProductEan, int64](dbProductEan, func(m model.ProductEan) int64 {
			return m.Ean
		})

		toCreate := list.FindMissing[int64](currentEans, previousEans)
		toDelete := list.FindMissing[int64](previousEans, currentEans)

		if err = s.productMetaRepository.CreateEANs(productId, toCreate, tx); err != nil {
			return err
		}

		if err = s.productMetaRepository.DeleteEANs(productId, toDelete, tx); err != nil {
			return err
		}

	} else if err = s.productMetaRepository.CreateEANs(productId, currentEans, tx); err != nil {
		// no records, create
		return err
	}

	return nil
}

func filterEANs(storeProduct datasource.StoreProduct) []int64 {
	var validEans []int64
	for _, ean := range storeProduct.EAN {
		if eanInt, err := strconv.Atoi(ean); err == nil {
			validEans = append(validEans, int64(eanInt))
		}
	}

	return validEans
}
