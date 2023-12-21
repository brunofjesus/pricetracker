package store

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/internal/repository"
)

const StoreTableName = "store"

type Store struct {
	StoreId int64  `db:"store_id"`
	Slug    string `db:"slug"`
	Name    string `db:"name"`
	Website string `db:"website"`
	Active  bool   `db:"active"`
}

type Repository struct {
	db *sql.DB
	qb *squirrel.StatementBuilderType
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
		qb: repository.QueryBuilder(db),
	}
}

func (r *Repository) FindStoreBySlug(slug string, tx *sql.Tx) (*Store, error) {
	qb := r.qb
	if tx != nil {
		qb = repository.QueryBuilder(tx)
	}

	q := qb.Select("store_id", "slug", "name", "website", "active").
		From(StoreTableName).
		Where(squirrel.Eq{"slug": slug})

	var store Store
	err := q.QueryRow().Scan(
		&store.StoreId,
		&store.Slug,
		&store.Name,
		&store.Website,
		&store.Active,
	)

	return &store, err
}

func (r *Repository) CreateStore(slug string, name string, website string, tx *sql.Tx) (int64, error) {
	qb := r.qb
	if tx != nil {
		qb = repository.QueryBuilder(tx)
	}

	q := qb.Insert(StoreTableName).
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
