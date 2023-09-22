package receiver

import (
	"log"
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/integration"
	"github.com/brunofjesus/pricetracker/catalog/src/service/product/creator"
	"github.com/brunofjesus/pricetracker/catalog/src/service/product/finder"
	"github.com/brunofjesus/pricetracker/catalog/src/service/product/updater"
)

var once sync.Once
var instance ProductReceiver

type ProductReceiver interface {
	Receive(storeProduct integration.StoreProduct)
}

type productReceiver struct {
	productFinder  finder.ProductFinder
	productCreator creator.ProductCreator
	productUpdater updater.ProductUpdater
}

func GetProductReceiver() ProductReceiver {
	once.Do(func() {
		instance = &productReceiver{
			productFinder:  finder.GetProductFinder(),
			productCreator: creator.GetProductCreator(),
			productUpdater: updater.GetProductUpdater(),
		}
	})
	return instance
}

// Receive implements Receiver.
func (s *productReceiver) Receive(storeProduct integration.StoreProduct) {
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
