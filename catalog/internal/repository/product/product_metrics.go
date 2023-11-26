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
	ProductId  int64  `db:"product_id" json:"product_id"`
	StoreId    int64  `db:"store_id" json:"store_id"`
	Name       string `db:"name" json:"name"`
	Brand      string `db:"brand" json:"brand"`
	Price      int    `db:"price" json:"price"`
	Available  bool   `db:"available" json:"available"`
	ImageUrl   string `db:"image_url" json:"image_url"`
	ProductUrl string `db:"product_url" json:"product_url"`

	Difference       decimal.Decimal `db:"diff" json:"diff"`
	DiscountPercent  decimal.Decimal `db:"discount_percent" json:"discount_percent"`
	Average          decimal.Decimal `db:"average" json:"average"`
	Minimum          decimal.Decimal `db:"minimum" json:"minimum"`
	Maximum          decimal.Decimal `db:"maximum" json:"maximum"`
	MetricEntryCount decimal.Decimal `db:"entries" json:"entries"`
	MetricDataSince  time.Time       `db:"metrics_since" json:"since"`
}

type ProductMetricsFilter struct {
	ProductId []int64

	StoreId    int
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
	MinMinimumPrice    float64
	MaxMinimumPrice    float64
	MinMaximumPrice    float64
	MaxMaximumPrice    float64
}

type ProductMetricsRepository interface {
	FindProductById(productId int64, tx *sql.Tx) (*ProductWithMetrics, error)
	FindProducts(offset int64, limit int, orderBy, direction string, filters *ProductMetricsFilter, tx *sql.Tx) ([]ProductWithMetrics, error)
	CountProducts(filters *ProductMetricsFilter, tx *sql.Tx) (int64, error)
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
func (r *productMetricsRepository) FindProducts(offset int64, limit int, orderBy, direction string, filters *ProductMetricsFilter, tx *sql.Tx) ([]ProductWithMetrics, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select(
		"product_id", "store_id", "name", "brand", "price", "available", "image_url", "product_url",
		"diff", "discount_percent", "average", "maximum", "minimum", "entries", "metrics_since",
	).From(ProductWithMetricsViewName)

	if filters != nil {
		q = appendFiltersToQuery(q, *filters)
	}

	q = q.OrderBy(fmt.Sprintf("%s %s", orderBy, direction)).
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
func (r *productMetricsRepository) CountProducts(filters *ProductMetricsFilter, tx *sql.Tx) (int64, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select("COUNT(*)").
		From(ProductWithMetricsViewName)

	if filters != nil {
		q = appendFiltersToQuery(q, *filters)
	}

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

func appendFiltersToQuery(q squirrel.SelectBuilder, filters ProductMetricsFilter) squirrel.SelectBuilder {
	f := squirrel.And{}

	if len(filters.ProductId) > 0 {
		f = append(f, squirrel.Eq{"product_id": filters.ProductId})
	}

	if filters.StoreId > -1 {
		f = append(f, squirrel.Eq{"store_id": filters.StoreId})
	}

	if len(filters.NameLike) > 0 {
		f = append(f, squirrel.Like{"name": "%" + filters.NameLike + "%"})
	}

	if len(filters.BrandLike) > 0 {
		f = append(f, squirrel.Like{"brand": "%" + filters.BrandLike + "%"})
	}

	if !nulltype.IsUndefined(filters.Available) {
		f = append(f, squirrel.Eq{"available": nulltype.IsTrue(filters.Available)})
	}

	if len(filters.ProductUrl) > 0 {
		f = append(f, squirrel.Eq{"product_url": filters.ProductUrl})
	}

	f = append(f, generateBetween("price", filters.MinPrice, filters.MaxPrice)...)
	f = append(f, generateBetween("diff", filters.MinDifference, filters.MaxDifference)...)
	f = append(f, generateBetween("discount_percent", filters.MinDiscountPercent, filters.MaxDiscountPercent)...)
	f = append(f, generateBetween("average", filters.MinAveragePrice, filters.MaxAveragePrice)...)
	f = append(f, generateBetween("maximum", filters.MinMaximumPrice, filters.MaxMaximumPrice)...)
	f = append(f, generateBetween("minimum", filters.MinMinimumPrice, filters.MaxMinimumPrice)...)

	q = q.Where(f)

	return q
}

func generateBetween(col string, minValue float64, maxValue float64) []squirrel.Sqlizer {
	var result []squirrel.Sqlizer = make([]squirrel.Sqlizer, 0, 2)

	if minValue > -1 {
		result = append(result, squirrel.GtOrEq{col: minValue})
	}
	if maxValue > -1 {
		result = append(result, squirrel.LtOrEq{col: maxValue})
	}

	return result
}
