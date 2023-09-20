package meta

import (
	"database/sql"
	"log"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/src/model"
	"github.com/brunofjesus/pricetracker/catalog/src/repository"
)

var once sync.Once
var instance ProductMetaRepository

type ProductMetaRepository interface {
	FindProductIdBySKU(sku []string, tx *sql.Tx) (int64, error)
	FindProductIdByEAN(ean []int64, tx *sql.Tx) (int64, error)
	CreateSKUs(productId int64, skus []string, tx *sql.Tx) error
	DeleteSKUs(productId int64, skus []string, tx *sql.Tx) error
	CreateEANs(productId int64, eans []int64, tx *sql.Tx) error
	DeleteEANs(productId int64, eans []int64, tx *sql.Tx) error
}

type productMetaRepository struct {
	db *sql.DB
	qb *squirrel.StatementBuilderType
}

func GetProductMetaRepository() ProductMetaRepository {
	once.Do(func() {
		db := repository.GetDatabaseConnection()

		instance = &productMetaRepository{
			db: db,
			qb: repository.QueryBuilder(db),
		}
	})
	return instance
}

// FindProductIdBySKU implements ProductMetaRepository.
func (r *productMetaRepository) FindProductIdBySKU(sku []string, tx *sql.Tx) (int64, error) {
	return r.findOne(tx, model.ProductSkuTableName, squirrel.Eq{"sku": sku})
}

// FindProductIdByEAN implements ProductMetaRepository.
func (r *productMetaRepository) FindProductIdByEAN(ean []int64, tx *sql.Tx) (int64, error) {
	return r.findOne(tx, model.ProductEanTableName, squirrel.Eq{"ean": ean})
}

// CreateEANs implements ProductMetaRepository.
func (r *productMetaRepository) CreateEANs(productId int64, eans []int64, tx *sql.Tx) error {
	var transaction *sql.Tx = tx
	var err error

	if transaction == nil {
		transaction, err = r.db.Begin()
		if err != nil {
			return err
		}
	}

	defer func() {
		if err != nil {
			log.Println(err)
			transaction.Rollback()
		}
	}()

	for _, ean := range eans {
		if err = r.createEan(productId, ean, transaction); err != nil {
			return err
		}
	}

	// only commit if the transaction was created by this method
	if tx == nil {
		err = transaction.Commit()
	}
	return err
}

// CreateSKUs implements ProductMetaRepository.
func (r *productMetaRepository) CreateSKUs(productId int64, skus []string, tx *sql.Tx) error {
	var transaction *sql.Tx = tx
	var err error

	if transaction == nil {
		transaction, err = r.db.Begin()
		if err != nil {
			return err
		}
	}

	defer func() {
		if err != nil {
			log.Println(err)
			transaction.Rollback()
		}
	}()

	for _, sku := range skus {
		if err = r.createSku(productId, sku, tx); err != nil {
			return err
		}
	}

	// only commit if the transaction was created by this method
	if tx == nil {
		err = transaction.Commit()
	}
	return err
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

func (r *productMetaRepository) findOne(tx *sql.Tx, tableName string, where any, args ...any) (int64, error) {
	qb := r.qb
	if tx != nil {
		qb = repository.QueryBuilder(tx)
	}

	q := qb.Select("product_id").
		From(tableName).
		Where(where, args...)

	var productId int64

	if err := q.QueryRow().Scan(&productId); err != nil {
		return -1, err
	}

	return productId, nil
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
