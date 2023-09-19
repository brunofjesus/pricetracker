package store

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"
	"github.com/brunofjesus/pricetracker/catalog/src/repository/store"
)

var once sync.Once
var instance StoreEnroller

type StoreEnroller interface {
	Enroll(store datasource.Store) error
}

type enroller struct {
	storeRepository store.StoreRepository
}

func GetStoreEnroller() StoreEnroller {
	once.Do(func() {
		instance = &enroller{
			storeRepository: store.GetStoreRepository(),
		}
	})
	return instance
}

// Enroll implements StoreEnroller.
func (s *enroller) Enroll(store datasource.Store) error {
	_, err := s.storeRepository.FindStoreBySlug(store.Slug(), nil)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		_, err = s.storeRepository.CreateStore(store.Slug(), store.Name(), store.Website(), nil)
	}

	return err
}
