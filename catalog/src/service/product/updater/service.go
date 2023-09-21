package updater

import (
	"database/sql"
	"sync"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"
	"github.com/brunofjesus/pricetracker/catalog/src/repository"

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

	// TODO: update SKUs and EANs

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

	return nil
}
