package store

import (
	"database/sql"
	"errors"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"
	"github.com/brunofjesus/pricetracker/catalog/src/repository/store"
)

type StoreEnroller interface {
	Enroll(store datasource.Store) error
}

type enroller struct {
	storeRepository store.StoreRepository
}

func NewStoreEnroller(storeRepository store.StoreRepository) StoreEnroller {
	return &enroller{
		storeRepository: storeRepository,
	}
}

// Enroll implements StoreEnroller.
func (s *enroller) Enroll(store datasource.Store) error {
	_, err := s.storeRepository.FindStoreBySlug(store.Slug(), nil)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		_, err = s.storeRepository.CreateStore(store.Slug(), store.Name(), store.Website(), nil)
	}

	return err
}
