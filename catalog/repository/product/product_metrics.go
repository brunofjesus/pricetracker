package product

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/model"
	"github.com/brunofjesus/pricetracker/catalog/repository"
)

var productMetricsOnce sync.Once
var productMetricsInstance ProductMetricsRepository

type ProductMetricsRepository interface {
	FindProductById(productId int64, tx *sql.Tx) (*model.ProductWithMetrics, error)
	FindProducts(offset int64, limit int, orderBy, direction string, tx *sql.Tx) ([]model.ProductWithMetrics, error)
	CountProducts(tx *sql.Tx) (int64, error)
}

type productMetricsRepository struct {
	db *sql.DB
	qb *squirrel.StatementBuilderType
}

func GetProductMetricsRepository() ProductMetricsRepository {
	productMetricsOnce.Do(func() {
		db := repository.GetDatabaseConnection()
		productMetricsInstance = &productMetricsRepository{
			db: db,
			qb: repository.QueryBuilder(db),
		}
	})

	return productMetricsInstance
}

// FindProductById implements ProductMetricsRepository.
func (r *productMetricsRepository) FindProductById(productId int64, tx *sql.Tx) (*model.ProductWithMetrics, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select(
		"product_id", "store_id", "name", "brand", "price", "available", "image_url", "product_url",
		"discount", "discount_percent", "average", "maximum", "minimum", "entries", "entries_since",
	).
		From(model.ProductWithMetricsViewName).
		Where(squirrel.Eq{"product_id": productId})

	var product model.ProductWithMetrics

	err := r.scanFullRow(q.QueryRow(), &product)

	return &product, err
}

// FindProducts implements ProductMetricsRepository.
func (r *productMetricsRepository) FindProducts(offset int64, limit int, orderBy string, direction string, tx *sql.Tx) ([]model.ProductWithMetrics, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select(
		"product_id", "store_id", "name", "brand", "price", "available", "image_url", "product_url",
		"discount", "discount_percent", "average", "maximum", "minimum", "entries", "entries_since",
	).
		From(model.ProductWithMetricsViewName).
		OrderBy(fmt.Sprintf("%s %s", orderBy, direction)).
		Offset(uint64(offset)).
		Limit(uint64(limit))

	var products []model.ProductWithMetrics
	rows, err := q.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product model.ProductWithMetrics

		if err := r.scanFullRow(rows, &product); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// CountProducts implements ProductMetricsRepository.
func (r *productMetricsRepository) CountProducts(tx *sql.Tx) (int64, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select("COUNT(*)").
		From(model.ProductWithMetricsViewName)

	var count int64
	err := q.QueryRow().Scan(&count)
	return count, err
}

func (r *productMetricsRepository) scanFullRow(row squirrel.RowScanner, product *model.ProductWithMetrics) error {
	return row.Scan(
		product.ProductId,
		product.StoreId,
		product.Name,
		product.Brand,
		product.Price,
		product.Available,
		product.ImageUrl,
		product.ProductUrl,

		product.Difference,
		product.DiscountPercent,
		product.Average,
		product.Maximum,
		product.Minimum,
		product.MetricEntryCount,
		product.MetricDataSince,
	)
}
