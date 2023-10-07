package product

import (
	"log"
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/integration"
)

var handlerOnce sync.Once
var handlerInstance ProductHandler

type ProductHandler interface {
	Handle(storeProduct integration.StoreProduct) error
}

type productHandler struct {
	productMatcher ProductMatcher
	productCreator ProductCreator
	productUpdater ProductUpdater
}

func GetProductHandler() ProductHandler {
	handlerOnce.Do(func() {
		handlerInstance = &productHandler{
			productMatcher: GetProductMatcher(),
			productCreator: GetProductCreator(),
			productUpdater: GetProductUpdater(),
		}
	})
	return handlerInstance
}

func (s *productHandler) Handle(storeProduct integration.StoreProduct) error {
	var err error
	if productId := s.productMatcher.Match(storeProduct); productId > 0 {
		err = s.productUpdater.Update(productId, storeProduct)
	} else {
		err = s.productCreator.Create(storeProduct)
	}

	if err != nil {
		log.Printf("error on receiver handler: %v on %v", err, storeProduct)
	}
	return err
}
