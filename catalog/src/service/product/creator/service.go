package creator

import (
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"

	price_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product"
	product_meta_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product/meta"
	store_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/store"
)

var once sync.Once
var instance ProductCreator

type ProductCreator interface {
	Create(storeProduct datasource.StoreProduct) error
}

type productCreator struct {
	storeRepository       store_repository.StoreRepository
	productRepository     product_repository.ProductRepository
	productMetaRepository product_meta_repository.ProductMetaRepository
	priceRepository       price_repository.PriceRepository
}

func GetProductCreator() ProductCreator {
	once.Do(func() {
		instance = &productCreator{
			storeRepository:       store_repository.GetStoreRepository(),
			productRepository:     product_repository.GetProductRepository(),
			productMetaRepository: product_meta_repository.GetProductMetaRepository(),
			priceRepository:       price_repository.GetPriceRepository(),
		}
	})
	return instance
}

// Create implements ProductUpdater.
func (*productCreator) Create(storeProduct datasource.StoreProduct) error {
	panic("unimplemented")
}
