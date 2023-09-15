package pingodoce

import (
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
	panic("unimplemented")
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
