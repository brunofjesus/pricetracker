package product

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/model"
	"github.com/brunofjesus/pricetracker/catalog/repository"
)

var productMetaOnce sync.Once
var productMetaInstance ProductMetaRepository

type ProductMetaRepository interface {
	FindProductIdBySKU(sku []string, storeSlug string, tx *sql.Tx) (int64, error)
	FindProductIdByEAN(ean []int64, storeSlug string, tx *sql.Tx) (int64, error)
	CreateSKUs(productId int64, skus []string, tx *sql.Tx) error
	DeleteSKUs(productId int64, skus []string, tx *sql.Tx) error
	CreateEANs(productId int64, eans []int64, tx *sql.Tx) error
	DeleteEANs(productId int64, eans []int64, tx *sql.Tx) error
	GetProductSKUs(productId int64, tx *sql.Tx) ([]model.ProductSku, error)
	GetProductEANs(productId int64, tx *sql.Tx) ([]model.ProductEan, error)
}

type productMetaRepository struct {
	db *sql.DB
	qb *squirrel.StatementBuilderType
}

func GetProductMetaRepository() ProductMetaRepository {
	productMetaOnce.Do(func() {
		db := repository.GetDatabaseConnection()

		productMetaInstance = &productMetaRepository{
			db: db,
			qb: repository.QueryBuilder(db),
		}
	})
	return productMetaInstance
}

// FindProductIdBySKU implements ProductMetaRepository.
func (r *productMetaRepository) FindProductIdBySKU(sku []string, storeSlug string, tx *sql.Tx) (int64, error) {
	return r.findOne(
		tx,
		model.ProductSkuTableName,
		squirrel.Eq{
			"slug": storeSlug,
			"sku":  sku,
		},
	)
}

// FindProductIdByEAN implements ProductMetaRepository.
func (r *productMetaRepository) FindProductIdByEAN(ean []int64, storeSlug string, tx *sql.Tx) (int64, error) {
	return r.findOne(
		tx,
		model.ProductEanTableName,
		squirrel.Eq{
			"slug": storeSlug,
			"ean":  ean,
		},
	)
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

// GetProductSKUs implements ProductMetaRepository.
func (r *productMetaRepository) GetProductSKUs(productId int64, tx *sql.Tx) ([]model.ProductSku, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select("product_id", "sku").
		From(model.ProductSkuTableName).
		Where(squirrel.Eq{"product_id": productId})

	var skus []model.ProductSku
	rows, err := q.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var sku model.ProductSku
		err := rows.Scan(
			&sku.ProductId,
			&sku.Sku,
		)

		if err != nil {
			return nil, err
		}

		skus = append(skus, sku)
	}

	return skus, nil
}

// GetProductEANs implements ProductMetaRepository.
func (r *productMetaRepository) GetProductEANs(productId int64, tx *sql.Tx) ([]model.ProductEan, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select("product_id", "ean").
		From(model.ProductEanTableName).
		Where(squirrel.Eq{"product_id": productId})

	var eans []model.ProductEan
	rows, err := q.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var ean model.ProductEan
		err := rows.Scan(
			&ean.ProductId,
			&ean.Ean,
		)

		if err != nil {
			return nil, err
		}

		eans = append(eans, ean)
	}

	return eans, nil
}

func (r *productMetaRepository) findOne(tx *sql.Tx, tableName string, where any, args ...any) (int64, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select("product_id").
		InnerJoin(fmt.Sprintf("%s USING (product_id)", model.ProductTableName)).
		InnerJoin(fmt.Sprintf("%s USING (store_id)", model.StoreTableName)).
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
