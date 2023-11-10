package store

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/model"
	"github.com/brunofjesus/pricetracker/catalog/repository/store"
)

var once sync.Once
var instance StoreHandler

type StoreHandler interface {
	Handle(store model.MqStore) error
}

type storeHandler struct {
	storeRepository store.StoreRepository
}

func GetStoreHandler() StoreHandler {
	once.Do(func() {
		instance = &storeHandler{
			storeRepository: store.GetStoreRepository(),
		}
	})
	return instance
}

// Handle implements StoreEnroller.
func (s *storeHandler) Handle(store model.MqStore) error {
	_, err := s.storeRepository.FindStoreBySlug(store.Slug, nil)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		_, err = s.storeRepository.CreateStore(store.Slug, store.Name, store.Website, nil)
	}

	return err
}
