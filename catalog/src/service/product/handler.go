package product

import (
	"log"
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/integration"
)

var handlerOnce sync.Once
var handlerInstance ProductHandler

type ProductHandler interface {
	Handle(storeProduct integration.StoreProduct)
}

type productHandler struct {
	productFinder  ProductFinder
	productCreator ProductCreator
	productUpdater ProductUpdater
}

func GetProductHandler() ProductHandler {
	handlerOnce.Do(func() {
		handlerInstance = &productHandler{
			productFinder:  GetProductFinder(),
			productCreator: GetProductCreator(),
			productUpdater: GetProductUpdater(),
		}
	})
	return handlerInstance
}

// Handle implements Receiver.
func (s *productHandler) Handle(storeProduct integration.StoreProduct) {
	var err error
	if productId := s.productFinder.Find(storeProduct); productId > 0 {
		err = s.productUpdater.Update(productId, storeProduct)
	} else {
		err = s.productCreator.Create(storeProduct)
	}

	if err != nil {
		log.Printf("error on receiver handler: %v", err)
	}
}
