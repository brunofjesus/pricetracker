package product

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/internal/repository"
	"github.com/brunofjesus/pricetracker/catalog/util/nulltype"
	"github.com/shopspring/decimal"
)

var productMetricsOnce sync.Once
var productMetricsInstance ProductMetricsRepository

const ProductWithMetricsViewName = "product_metrics"

type ProductWithMetrics struct {
	ProductId  int64  `db:"product_id"`
	StoreId    int64  `db:"store_id"`
	Name       string `db:"name"`
	Brand      string `db:"brand"`
	Price      int    `db:"price"`
	Available  bool   `db:"available"`
	ImageUrl   string `db:"image_url"`
	ProductUrl string `db:"product_url"`

	Difference       decimal.Decimal `db:"diff"`
	DiscountPercent  decimal.Decimal `db:"discount_percent"`
	Average          decimal.Decimal `db:"average"`
	Maximum          decimal.Decimal `db:"maximum"`
	Minimum          decimal.Decimal `db:"minimum"`
	MetricEntryCount decimal.Decimal `db:"entries"`
	MetricDataSince  time.Time       `db:"metrics_since"`
}

// TODO: use in query
type ProductMetricsFilter struct {
	MinPrice   float64
	MaxPrice   float64
	NameLike   string
	BrandLike  string
	Available  nulltype.NullBoolean
	ProductUrl string

	MinDifference      float64
	MaxDifference      float64
	MinDiscountPercent float64
	MaxDiscountPercent float64
	MinAveragePrice    float64
	MaxAveragePrice    float64
	MinMaximumPrice    float64
	MaxMaximumPrice    float64
	MinMinimumPrice    float64
	MaxMinimumPrice    float64
}

type ProductMetricsRepository interface {
	FindProductById(productId int64, tx *sql.Tx) (*ProductWithMetrics, error)
	FindProducts(offset int64, limit int, orderBy, direction string, tx *sql.Tx) ([]ProductWithMetrics, error)
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
func (r *productMetricsRepository) FindProductById(productId int64, tx *sql.Tx) (*ProductWithMetrics, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select(
		"product_id", "store_id", "name", "brand", "price", "available", "image_url", "product_url",
		"diff", "discount_percent", "average", "maximum", "minimum", "entries", "metrics_since",
	).
		From(ProductWithMetricsViewName).
		Where(squirrel.Eq{"product_id": productId})

	var product ProductWithMetrics

	err := r.scanFullRow(q.QueryRow(), &product)

	return &product, err
}

// FindProducts implements ProductMetricsRepository.
func (r *productMetricsRepository) FindProducts(offset int64, limit int, orderBy string, direction string, tx *sql.Tx) ([]ProductWithMetrics, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select(
		"product_id", "store_id", "name", "brand", "price", "available", "image_url", "product_url",
		"diff", "discount_percent", "average", "maximum", "minimum", "entries", "metrics_since",
	).
		From(ProductWithMetricsViewName).
		OrderBy(fmt.Sprintf("%s %s", orderBy, direction)).
		Offset(uint64(offset)).
		Limit(uint64(limit))

	var products []ProductWithMetrics
	rows, err := q.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product ProductWithMetrics

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
		From(ProductWithMetricsViewName)

	var count int64
	err := q.QueryRow().Scan(&count)
	return count, err
}

func (r *productMetricsRepository) scanFullRow(row squirrel.RowScanner, product *ProductWithMetrics) error {
	return row.Scan(
		&product.ProductId,
		&product.StoreId,
		&product.Name,
		&product.Brand,
		&product.Price,
		&product.Available,
		&product.ImageUrl,
		&product.ProductUrl,

		&product.Difference,
		&product.DiscountPercent,
		&product.Average,
		&product.Maximum,
		&product.Minimum,
		&product.MetricEntryCount,
		&product.MetricDataSince,
	)
}
