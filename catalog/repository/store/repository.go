package store

import (
	"database/sql"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/model"
	"github.com/brunofjesus/pricetracker/catalog/repository"
)

var once sync.Once
var instance StoreRepository

type StoreRepository interface {
	FindStoreBySlug(slug string, tx *sql.Tx) (*model.Store, error)
	CreateStore(slug, name, website string, tx *sql.Tx) (int64, error)
}

type storeRepository struct {
	db *sql.DB
	qb *squirrel.StatementBuilderType
}

func GetStoreRepository() StoreRepository {
	once.Do(func() {
		db := repository.GetDatabaseConnection()

		instance = &storeRepository{
			db: db,
			qb: repository.QueryBuilder(db),
		}
	})

	return instance
}

func (r *storeRepository) FindStoreBySlug(slug string, tx *sql.Tx) (*model.Store, error) {
	qb := r.qb
	if tx != nil {
		qb = repository.QueryBuilder(tx)
	}

	q := qb.Select("store_id", "slug", "name", "website", "active").
		From(model.StoreTableName).
		Where(squirrel.Eq{"slug": slug})

	var store model.Store
	err := q.QueryRow().Scan(
		&store.StoreId,
		&store.Slug,
		&store.Name,
		&store.Website,
		&store.Active,
	)

	return &store, err
}

func (r *storeRepository) CreateStore(slug string, name string, website string, tx *sql.Tx) (int64, error) {
	qb := r.qb
	if tx != nil {
		qb = repository.QueryBuilder(tx)
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
