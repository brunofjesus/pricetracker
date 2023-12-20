package store

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/internal/repository/store"
)

var once sync.Once
var instance StoreHandler

type StoreHandler interface {
	Handle(ctx context.Context, store MqStore) error
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

func (s *storeHandler) Handle(ctx context.Context, store MqStore) error {
	_, err := s.storeRepository.FindStoreBySlug(store.Slug, nil)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		_, err = s.storeRepository.CreateStore(store.Slug, store.Name, store.Website, nil)
	}

	return err
}
