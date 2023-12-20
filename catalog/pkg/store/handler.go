package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/brunofjesus/pricetracker/catalog/internal/repository/store"
)

type Handler struct {
	Repository *store.Repository
}

func (s *Handler) Handle(ctx context.Context, store MqStore) error {
	_, err := s.Repository.FindStoreBySlug(store.Slug, nil)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		_, err = s.Repository.CreateStore(store.Slug, store.Name, store.Website, nil)
	}

	return err
}
