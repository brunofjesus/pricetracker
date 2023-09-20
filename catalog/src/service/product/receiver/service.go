package receiver

import (
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"
	"github.com/brunofjesus/pricetracker/catalog/src/service/product/creator"
	"github.com/brunofjesus/pricetracker/catalog/src/service/product/finder"
	"github.com/brunofjesus/pricetracker/catalog/src/service/product/updater"
)

var once sync.Once
var instance ProductReceiver

type ProductReceiver interface {
	Receive(storeProduct datasource.StoreProduct)
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
func (s *productReceiver) Receive(storeProduct datasource.StoreProduct) {
	if productId := s.productFinder.Find(storeProduct); productId > 0 {
		s.productUpdater.Update(storeProduct)
	} else {
		s.productCreator.Create(storeProduct)
	}
}
