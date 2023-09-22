package main

import (
	"sync"
)

var once sync.Once
var instance PingoDoceStore

type PingoDoceStore struct {
}

func Instance() Store {
	once.Do(
		func() {
			instance = PingoDoceStore{}
		},
	)

	return &instance
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
