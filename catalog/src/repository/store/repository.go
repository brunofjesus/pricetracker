package store

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/src/model"
	"github.com/brunofjesus/pricetracker/catalog/src/repository"
)

type StoreRepository interface {
	CreateStore(slug, name, website string, transaction *sql.Tx) (int64, error)
}

type storeRepository struct {
	db *sql.DB
	qb *squirrel.StatementBuilderType
}

func NewStoreRepository(db *sql.DB, qb *squirrel.StatementBuilderType) StoreRepository {
	return &storeRepository{
		db: db,
		qb: qb,
	}
}

func (r *storeRepository) CreateStore(slug string, name string, website string, transaction *sql.Tx) (int64, error) {
	qb := r.qb
	if transaction != nil {
		qb = repository.QueryBuilder(transaction)
	}

	q := qb.Insert(model.StoreTableName).
		Columns("slug", "name", "website", "active").
		Values(slug, name, website, true).
		Suffix("RETURNING store_id")

	var id int64
	err := q.QueryRow().Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}
