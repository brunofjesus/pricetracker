package meta

import (
	"database/sql"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/src/model"
	"github.com/brunofjesus/pricetracker/catalog/src/repository"
)

type ProductMetaRepository interface {
	CreateSKUs(productId int64, skus []string, tx *sql.Tx) error
	DeleteSKUs(productId int64, skus []string, tx *sql.Tx) error
	CreateEANs(productId int64, eans []int64, tx *sql.Tx) error
	DeleteEANs(productId int64, eans []int64, tx *sql.Tx) error
}

type productMetaRepository struct {
	db *sql.DB
	qb *squirrel.StatementBuilderType
}

func NewProductMetaRepository(db *sql.DB) ProductMetaRepository {
	return &productMetaRepository{
		db: db,
		qb: repository.QueryBuilder(db),
	}
}

// CreateEANs implements ProductMetaRepository.
func (r *productMetaRepository) CreateEANs(productId int64, eans []int64, tx *sql.Tx) error {
	var err error

	if tx == nil {
		tx, err = r.db.Begin()
		if err != nil {
			return err
		}
	}

	defer func() {
		if err != nil {
			log.Println(err)
			tx.Rollback()
		}
	}()

	for _, ean := range eans {
		if err = r.createEan(productId, ean, tx); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// CreateSKUs implements ProductMetaRepository.
func (r *productMetaRepository) CreateSKUs(productId int64, skus []string, tx *sql.Tx) error {
	var err error

	if tx == nil {
		tx, err = r.db.Begin()
		if err != nil {
			return err
		}
	}

	defer func() {
		if err != nil {
			log.Println(err)
			tx.Rollback()
		}
	}()

	for _, sku := range skus {
		if err = r.createSku(productId, sku, tx); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteEANs implements ProductMetaRepository.
func (r *productMetaRepository) DeleteEANs(productId int64, eans []int64, tx *sql.Tx) error {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Delete(model.ProductEanTableName).
		Where(
			squirrel.And{
				squirrel.Eq{
					"product_id": productId,
				},
				squirrel.Eq{
					"ean": eans,
				},
			},
		)

	_, err := q.Exec()
	if err != nil {
		return err
	}
	return nil
}

// DeleteSKUs implements ProductMetaRepository.
func (r *productMetaRepository) DeleteSKUs(productId int64, skus []string, tx *sql.Tx) error {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Delete(model.ProductSkuTableName).
		Where(
			squirrel.And{
				squirrel.Eq{
					"product_id": productId,
				},
				squirrel.Eq{
					"sku": skus,
				},
			},
		)

	_, err := q.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (r *productMetaRepository) createEan(productId int64, ean int64, tx *sql.Tx) error {
	qb := repository.QueryBuilder(tx)

	q := qb.Insert(model.ProductEanTableName).
		Columns("product_id", "ean").
		Values(productId, ean)

	_, err := q.Exec()

	return err
}

func (r *productMetaRepository) createSku(productId int64, sku string, tx *sql.Tx) error {
	qb := repository.QueryBuilder(tx)

	q := qb.Insert(model.ProductSkuTableName).
		Columns("product_id", "sku").
		Values(productId, sku)

	_, err := q.Exec()

	return err
}
