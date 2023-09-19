package receiver

import (
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"
	price_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product"
	product_meta_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product/meta"
	store_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/store"
)

var once sync.Once
var instance ProductReceiver

type ProductReceiver interface {
	Receive(storeProduct datasource.StoreProduct)
}

type productReceiver struct {
	storeRepository       store_repository.StoreRepository
	productRepository     product_repository.ProductRepository
	productMetaRepository product_meta_repository.ProductMetaRepository
	priceRepository       price_repository.PriceRepository
}

func GetProductReceiver() ProductReceiver {
	once.Do(func() {
		instance = &productReceiver{
			storeRepository:       store_repository.GetStoreRepository(),
			productRepository:     product_repository.GetProductRepository(),
			productMetaRepository: product_meta_repository.GetProductMetaRepository(),
			priceRepository:       price_repository.GetPriceRepository(),
		}
	})
	return instance
}

// Receive implements Receiver.
func (*productReceiver) Receive(storeProduct datasource.StoreProduct) {
	panic("unimplemented")
}
