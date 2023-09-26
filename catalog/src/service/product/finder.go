package product

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/integration"
	price_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product"
	product_meta_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product/meta"
	store_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/store"
)

var finderOnce sync.Once
var finderInstance ProductFinder

type ProductFinder interface {
	Find(storeProduct integration.StoreProduct) int64
}

type productFinder struct {
	storeRepository       store_repository.StoreRepository
	productRepository     product_repository.ProductRepository
	productMetaRepository product_meta_repository.ProductMetaRepository
	priceRepository       price_repository.PriceRepository
}

func GetProductFinder() ProductFinder {
	finderOnce.Do(func() {
		finderInstance = &productFinder{
			storeRepository:       store_repository.GetStoreRepository(),
			productRepository:     product_repository.GetProductRepository(),
			productMetaRepository: product_meta_repository.GetProductMetaRepository(),
			priceRepository:       price_repository.GetPriceRepository(),
		}
	})
	return finderInstance
}

// Create implements ProductUpdater.
func (s *productFinder) Find(storeProduct integration.StoreProduct) int64 {
	var searchFunctions []func(integration.StoreProduct) int64
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

func (s *productFinder) findByEan(storeProduct integration.StoreProduct) int64 {
	var validEans []int64
	for _, ean := range storeProduct.EAN {
		if eanInt, err := strconv.Atoi(ean); err == nil {
			validEans = append(validEans, int64(eanInt))
		}
	}

	if len(validEans) > 0 {
		productId, err := s.productMetaRepository.FindProductIdByEAN(validEans, storeProduct.StoreSlug, nil)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("error finding by ean %v: %v", validEans, err)
		} else if err == nil {
			return productId
		}
	}

	return -1
}

func (s *productFinder) findBySku(storeProduct integration.StoreProduct) int64 {
	productId, err := s.productMetaRepository.FindProductIdBySKU(storeProduct.SKU, storeProduct.StoreSlug, nil)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("error finding by sku %v: %v", storeProduct.SKU, err)
	} else if err == nil {
		return productId
	}

	return -1
}

func (s *productFinder) findByProductUrl(storeProduct integration.StoreProduct) int64 {
	product, err := s.productRepository.FindProductByUrl(storeProduct.Link, nil)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("error finding by url %v: %v", storeProduct.Link, err)
	} else if err == nil {
		return product.ProductId
	}

	return -1
}
