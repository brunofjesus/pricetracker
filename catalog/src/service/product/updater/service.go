package updater

import (
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"

	price_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/price"
	product_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product"
	product_meta_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/product/meta"
	store_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/store"
)

var once sync.Once
var instance ProductUpdater

type ProductUpdater interface {
	Update(storeProduct datasource.StoreProduct) error
}

type productUpdater struct {
	storeRepository       store_repository.StoreRepository
	productRepository     product_repository.ProductRepository
	productMetaRepository product_meta_repository.ProductMetaRepository
	priceRepository       price_repository.PriceRepository
}

func GetProductUpdater() ProductUpdater {
	once.Do(func() {
		instance = &productUpdater{
			storeRepository:       store_repository.GetStoreRepository(),
			productRepository:     product_repository.GetProductRepository(),
			productMetaRepository: product_meta_repository.GetProductMetaRepository(),
			priceRepository:       price_repository.GetPriceRepository(),
		}
	})
	return instance
}

// Update implements ProductUpdater.
func (*productUpdater) Update(storeProduct datasource.StoreProduct) error {
	panic("unimplemented")
}
