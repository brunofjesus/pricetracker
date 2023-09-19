package pingodoce

import (
	"log"
	"sync"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"
)

var once sync.Once
var instance PingoDoceStore

type PingoDoceStore struct {
}

func Instance() datasource.Store {
	once.Do(
		func() {
			instance = PingoDoceStore{}
		},
	)

	return &instance
}

// Crawl implements store.Store.
func (s *PingoDoceStore) Crawl(productChannel chan datasource.StoreProduct) {

	index := 0
	total := 100
	for index < total {
		response, err := search(index, 100)

		if err != nil {
			log.Default().Print(err)
			break
		}

		total = response.Sections.Null.Total
		index = index + 100

		for _, product := range response.Sections.Null.Products {
			storeProduct := mapPingoDoceProductToStoreProduct(product.Source)
			productChannel <- storeProduct
		}

		time.Sleep(1 * time.Second) // be polite
	}
}

// Slug implements store.Store.
func (s *PingoDoceStore) Slug() string {
	return "pingo-doce"
}

// Name implements store.Store.
func (s *PingoDoceStore) Name() string {
	return "Pingo Doce"
}

// Website implements store.Store.
func (*PingoDoceStore) Website() string {
	return "https://mercadao.pt/store/pingo-doce"
}
