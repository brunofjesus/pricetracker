package finder

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"

	price_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product"
	product_meta_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product/meta"
	store_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/store"
)

var once sync.Once
var instance ProductFinder

type ProductFinder interface {
	Find(storeProduct datasource.StoreProduct) int64
}

type productFinder struct {
	storeRepository       store_repository.StoreRepository
	productRepository     product_repository.ProductRepository
	productMetaRepository product_meta_repository.ProductMetaRepository
	priceRepository       price_repository.PriceRepository
}

func GetProductFinder() ProductFinder {
	once.Do(func() {
		instance = &productFinder{
			storeRepository:       store_repository.GetStoreRepository(),
			productRepository:     product_repository.GetProductRepository(),
			productMetaRepository: product_meta_repository.GetProductMetaRepository(),
			priceRepository:       price_repository.GetPriceRepository(),
		}
	})
	return instance
}

// Create implements ProductUpdater.
func (s *productFinder) Find(storeProduct datasource.StoreProduct) int64 {
	var searchFunctions []func(datasource.StoreProduct) int64
	searchFunctions = append(searchFunctions, s.findByEan)
	searchFunctions = append(searchFunctions, s.findBySku)
	searchFunctions = append(searchFunctions, s.findByProductUrl)

	for _, searchFunction := range searchFunctions {
		productId := searchFunction(storeProduct)

		if productId > 0 {
			return productId
		}
	}

	return -1
}

// FIXME: we need to consider the store slug as well
func (s *productFinder) findByEan(storeProduct datasource.StoreProduct) int64 {
	var validEans []int64
	for _, ean := range storeProduct.EAN {
		if eanInt, err := strconv.Atoi(ean); err == nil {
			validEans = append(validEans, int64(eanInt))
		}
	}

	if len(validEans) > 0 {
		productId, err := s.productMetaRepository.FindProductIdByEAN(validEans, nil)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("error finding by ean %v: %v", validEans, err)
		} else if err == nil {
			return productId
		}
	}

	return -1
}

// FIXME: we need to consider the store slug as well
func (s *productFinder) findBySku(storeProduct datasource.StoreProduct) int64 {
	productId, err := s.productMetaRepository.FindProductIdBySKU(storeProduct.SKU, nil)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("error finding by sku %v: %v", storeProduct.SKU, err)
	} else if err == nil {
		return productId
	}

	return -1
}

func (s *productFinder) findByProductUrl(storeProduct datasource.StoreProduct) int64 {
	product, err := s.productRepository.FindProductByUrl(storeProduct.Link, nil)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("error finding by url %v: %v", storeProduct.Link, err)
	} else if err == nil {
		return product.ProductId
	}

	return -1
}
