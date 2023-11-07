package product

import (
	"database/sql"
	"sync"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/model"
	"github.com/brunofjesus/pricetracker/catalog/src/repository"

	price_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product"
	store_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/store"
)

var creatorOnce sync.Once
var creatorInstance ProductCreator

type ProductCreator interface {
	Create(storeProduct model.MqStoreProduct) error
}

type productCreator struct {
	db                    *sql.DB
	storeRepository       store_repository.StoreRepository
	productRepository     product_repository.ProductRepository
	productMetaRepository product_repository.ProductMetaRepository
	priceRepository       price_repository.PriceRepository
}

func GetProductCreator() ProductCreator {
	creatorOnce.Do(func() {
		creatorInstance = &productCreator{
			db:                    repository.GetDatabaseConnection(),
			storeRepository:       store_repository.GetStoreRepository(),
			productRepository:     product_repository.GetProductRepository(),
			productMetaRepository: product_repository.GetProductMetaRepository(),
			priceRepository:       price_repository.GetPriceRepository(),
		}
	})
	return creatorInstance
}

// Create implements ProductUpdater.
func (s *productCreator) Create(storeProduct model.MqStoreProduct) error {
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
