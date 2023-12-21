package product

import (
	"database/sql"
	"time"

	price_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/product"
	store_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/store"
)

type Creator struct {
	DB                    *sql.DB
	StoreRepository       *store_repository.Repository
	ProductRepository     *product_repository.Repository
	ProductMetaRepository *product_repository.MetaRepository
	PriceRepository       *price_repository.Repository
}

func (s *Creator) Create(storeProduct MqStoreProduct) error {
	tx, err := s.DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	store, err := s.StoreRepository.FindStoreBySlug(storeProduct.StoreSlug, tx)
	if err != nil {
		return err
	}

	productId, err := s.ProductRepository.CreateProduct(
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
		s.ProductMetaRepository.CreateEANs(
			productId,
			filterNumbers(storeProduct.EAN),
			tx,
		)
	}

	if len(storeProduct.SKU) > 0 {
		s.ProductMetaRepository.CreateSKUs(
			productId,
			storeProduct.SKU,
			tx,
		)
	}

	err = s.PriceRepository.CreatePrice(productId, storeProduct.Price, time.Now(), tx)

	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
