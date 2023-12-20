package product

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	price_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/product"
	store_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/store"
)

type ProductMatcher struct {
	StoreRepository       store_repository.StoreRepository
	ProductRepository     product_repository.ProductRepository
	ProductMetaRepository product_repository.ProductMetaRepository
	PriceRepository       price_repository.PriceRepository
}

func (s *ProductMatcher) Match(storeProduct MqStoreProduct) int64 {
	var searchFunctions []func(MqStoreProduct) int64
	searchFunctions = append(searchFunctions, s.findByProductUrl)
	searchFunctions = append(searchFunctions, s.findByEan)
	searchFunctions = append(searchFunctions, s.findBySku)

	for _, searchFunction := range searchFunctions {
		productId := searchFunction(storeProduct)

		if productId > 0 {
			return productId
		}
	}

	return -1
}

func (s *ProductMatcher) findByEan(storeProduct MqStoreProduct) int64 {
	var validEans []int64
	for _, ean := range storeProduct.EAN {
		if eanInt, err := strconv.Atoi(ean); err == nil {
			validEans = append(validEans, int64(eanInt))
		}
	}

	if len(validEans) > 0 {
		productId, err := s.ProductMetaRepository.FindProductIdByEAN(validEans, storeProduct.StoreSlug, nil)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("error finding by ean %v: %v", validEans, err)
		} else if err == nil {
			return productId
		}
	}

	return -1
}

func (s *ProductMatcher) findBySku(storeProduct MqStoreProduct) int64 {
	productId, err := s.ProductMetaRepository.FindProductIdBySKU(storeProduct.SKU, storeProduct.StoreSlug, nil)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("error finding by sku %v: %v", storeProduct.SKU, err)
	} else if err == nil {
		return productId
	}

	return -1
}

func (s *ProductMatcher) findByProductUrl(storeProduct MqStoreProduct) int64 {
	product, err := s.ProductRepository.FindProductByUrl(storeProduct.Link, nil)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("error finding by url %v: %v", storeProduct.Link, err)
	} else if err == nil {
		return product.ProductId
	}

	return -1
}
