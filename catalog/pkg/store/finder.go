package store

import (
	store_repository "github.com/brunofjesus/pricetracker/catalog/internal/repository/store"
)

type Finder struct {
	StoreRepository *store_repository.Repository
}

func (s *Finder) FindStores() ([]store_repository.Store, error) {
	return s.StoreRepository.FindStores()
}

func (s *Finder) FindStoreBySlug(slug string) (*store_repository.Store, error) {
	return s.StoreRepository.FindStoreBySlug(slug, nil)
}
