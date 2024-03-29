package product

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/internal/repository"
	"github.com/brunofjesus/pricetracker/catalog/internal/repository/stats"
	"github.com/brunofjesus/pricetracker/catalog/internal/repository/store"
)

const ProductWithStatsViewName = "product_with_stats"

type ProductWithStatsFilter struct {
	ProductId []int64

	StoreId    sql.NullInt32
	MinPrice   sql.NullFloat64
	MaxPrice   sql.NullFloat64
	NameLike   sql.NullString
	BrandLike  sql.NullString
	Available  sql.NullBool
	ProductUrl sql.NullString

	MinDifference      sql.NullFloat64
	MaxDifference      sql.NullFloat64
	MinDiscountPercent sql.NullFloat64
	MaxDiscountPercent sql.NullFloat64
	MinAveragePrice    sql.NullFloat64
	MaxAveragePrice    sql.NullFloat64
	MinMinimumPrice    sql.NullFloat64
	MaxMinimumPrice    sql.NullFloat64
	MinMaximumPrice    sql.NullFloat64
	MaxMaximumPrice    sql.NullFloat64
}

var cols = []string{
	"product_id", "store_id", "store_name", "store_slug", "store_website", "name", "brand", "price", "currency",
	"available", "image_url", "product_url", "difference", "discount_percent", "average", "maximum", "minimum",
	"entries",
}

type ProductWithStatsRepository struct {
	db *sql.DB
	qb *squirrel.StatementBuilderType
}

func NewMetricsRepository(db *sql.DB) *ProductWithStatsRepository {
	return &ProductWithStatsRepository{
		db: db,
		qb: repository.QueryBuilder(db),
	}
}

func (r *ProductWithStatsRepository) QuickSearch(tsExpression string) ([]Product, error) {
	tsQuery := fmt.Sprintf("to_tsquery('english', '%s')", tsExpression)

	q := r.qb.Select(
		"product_id", "store_id", "store.name", "store.slug", "store.website",
		"product.name", "brand", "product.price", "currency", "available", "image_url", "product_url",
		"difference", "discount_percent", "average", "maximum", "minimum",
		"entries").
		From(ProductTableName).
		InnerJoin(fmt.Sprintf("%s USING (product_id)", stats.ProductStatsTableName)).
		InnerJoin(fmt.Sprintf("%s as store USING (store_id)", store.StoreTableName)).
		Where(fmt.Sprintf("search_vector @@ %s", tsQuery)).
		OrderBy(fmt.Sprintf("ts_rank(search_vector, %s)", tsQuery)).
		Limit(20).
		Offset(0)

	var products []Product
	rows, err := q.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product Product

		if err := r.scanFullRow(rows, &product); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductWithStatsRepository) FindProductById(productId int64, tx *sql.Tx) (*Product, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select(cols...).From(ProductWithStatsViewName).Where(squirrel.Eq{"product_id": productId})

	var product Product

	err := r.scanFullRow(q.QueryRow(), &product)

	return &product, err
}

func (r *ProductWithStatsRepository) FindProducts(offset int64, limit int, orderBy, direction string, filters *ProductWithStatsFilter, tx *sql.Tx) ([]Product, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select(cols...).From(ProductWithStatsViewName)
	if filters != nil {
		q = appendFiltersToQuery(q, *filters)
	}

	q = q.OrderBy(fmt.Sprintf("%s %s", orderBy, direction)).
		Offset(uint64(offset)).
		Limit(uint64(limit))

	var products []Product
	rows, err := q.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product Product

		if err := r.scanFullRow(rows, &product); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductWithStatsRepository) CountProducts(filters *ProductWithStatsFilter, tx *sql.Tx) (int64, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select("COUNT(*)").From(ProductWithStatsViewName)

	if filters != nil {
		q = appendFiltersToQuery(q, *filters)
	}

	var count int64
	err := q.QueryRow().Scan(&count)
	return count, err
}

func (r *ProductWithStatsRepository) scanFullRow(row squirrel.RowScanner, product *Product) error {
	store := Store{}
	stats := Statistics{}

	err := row.Scan(
		&product.ProductId,
		&store.StoreId,
		&store.Name,
		&store.Slug,
		&store.Website,
		&product.Name,
		&product.Brand,
		&product.Price,
		&product.Currency,
		&product.Available,
		&product.ImageUrl,
		&product.ProductUrl,

		&stats.Difference,
		&stats.DiscountPercent,
		&stats.Average,
		&stats.Maximum,
		&stats.Minimum,
		&stats.MetricEntryCount,
	)

	if err != nil {
		return err
	}

	product.Store = &store
	product.Statistics = &stats

	return nil
}

func appendFiltersToQuery(q squirrel.SelectBuilder, filters ProductWithStatsFilter) squirrel.SelectBuilder {
	f := squirrel.And{}

	if len(filters.ProductId) > 0 {
		f = append(f, squirrel.Eq{"product_id": filters.ProductId})
	}

	if filters.StoreId.Valid {
		f = append(f, squirrel.Eq{"store_id": filters.StoreId.Int32})
	}

	if filters.NameLike.Valid {
		f = append(f, squirrel.Like{"name": "%" + filters.NameLike.String + "%"})
	}

	if filters.BrandLike.Valid {
		f = append(f, squirrel.Like{"brand": "%" + filters.BrandLike.String + "%"})
	}

	if filters.Available.Valid {
		f = append(f, squirrel.Eq{"available": filters.Available.Bool})
	}

	if filters.ProductUrl.Valid {
		f = append(f, squirrel.Eq{"product_url": filters.ProductUrl})
	}

	f = append(f, generateBetween("price", filters.MinPrice, filters.MaxPrice)...)
	f = append(f, generateBetween("difference", filters.MinDifference, filters.MaxDifference)...)
	f = append(f, generateBetween("discount_percent", filters.MinDiscountPercent, filters.MaxDiscountPercent)...)
	f = append(f, generateBetween("average", filters.MinAveragePrice, filters.MaxAveragePrice)...)
	f = append(f, generateBetween("maximum", filters.MinMaximumPrice, filters.MaxMaximumPrice)...)
	f = append(f, generateBetween("minimum", filters.MinMinimumPrice, filters.MaxMinimumPrice)...)

	q = q.Where(f)

	return q
}

func generateBetween(col string, minValue sql.NullFloat64, maxValue sql.NullFloat64) []squirrel.Sqlizer {
	var result []squirrel.Sqlizer = make([]squirrel.Sqlizer, 0, 2)

	if minValue.Valid {
		result = append(result, squirrel.GtOrEq{col: minValue.Float64})
	}
	if maxValue.Valid {
		result = append(result, squirrel.LtOrEq{col: maxValue.Float64})
	}

	return result
}
