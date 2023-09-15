package pingodoce

import (
	"fmt"
	"log"
	"sync"

	"github.com/brunofjesus/pricetracker/store"
)

var once sync.Once
var instance PingoDoceStore

type PingoDoceStore struct {
}

func Instance() store.Store {
	once.Do(
		func() {
			instance = PingoDoceStore{}
		},
	)

	return &instance
}

// Crawl implements store.Store.
func (s *PingoDoceStore) Crawl(productChannel chan store.StoreProduct) {

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
			fmt.Printf("%s: %f\n", product.Source.ShortDescription, product.Source.UnitPrice)
		}
	}
}

// Id implements store.Store.
func (s *PingoDoceStore) Id() string {
	panic("unimplemented")
}

// Name implements store.Store.
func (s *PingoDoceStore) Name() string {
	panic("unimplemented")
}

// Website implements store.Store.
func (*PingoDoceStore) Website() string {
	panic("unimplemented")
}
