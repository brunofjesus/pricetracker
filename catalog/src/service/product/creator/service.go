package creator

import (
	"database/sql"
	"strconv"
	"sync"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/repository"

	"github.com/brunofjesus/pricetracker/catalog/src/integration"
	price_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product"
	product_meta_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product/meta"
	store_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/store"
)

var once sync.Once
var instance ProductCreator

type ProductCreator interface {
	Create(storeProduct integration.StoreProduct) error
}

type productCreator struct {
	db                    *sql.DB
	storeRepository       store_repository.StoreRepository
	productRepository     product_repository.ProductRepository
	productMetaRepository product_meta_repository.ProductMetaRepository
	priceRepository       price_repository.PriceRepository
}

func GetProductCreator() ProductCreator {
	once.Do(func() {
		instance = &productCreator{
			db:                    repository.GetDatabaseConnection(),
			storeRepository:       store_repository.GetStoreRepository(),
			productRepository:     product_repository.GetProductRepository(),
			productMetaRepository: product_meta_repository.GetProductMetaRepository(),
			priceRepository:       price_repository.GetPriceRepository(),
		}
	})
	return instance
}

// Create implements ProductUpdater.
func (s *productCreator) Create(storeProduct integration.StoreProduct) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	store, err := s.storeRepository.FindStoreBySlug(storeProduct.StoreSlug, tx)
	if err != nil {
		return err
	}

	productId, err := s.productRepository.CreateProduct(
		store.StoreId,
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

	if len(storeProduct.EAN) > 0 {
		s.productMetaRepository.CreateEANs(
			productId,
			filterEANs(storeProduct),
			tx,
		)
	}

	if len(storeProduct.SKU) > 0 {
		s.productMetaRepository.CreateSKUs(
			productId,
			storeProduct.SKU,
			tx,
		)
	}

	err = s.priceRepository.CreatePrice(productId, storeProduct.Price, time.Now(), tx)

	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func filterEANs(storeProduct integration.StoreProduct) []int64 {
	var validEans []int64
	for _, ean := range storeProduct.EAN {
		if eanInt, err := strconv.Atoi(ean); err == nil {
			validEans = append(validEans, int64(eanInt))
		}
	}

	return validEans
}
