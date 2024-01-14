package stats

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/internal/repository"
	"github.com/shopspring/decimal"
)

const ProductStatsTableName = "product_stats"

var productStatsTableColumns = []string{
	"product_id", "price", "difference", "discount_percent", "average", "minimum", "maximum", "entries",
}

type ProductStats struct {
	ProductId        int64           `db:"product_id" json:"productId"`
	Price            int             `db:"price" json:"price"`
	Difference       int             `db:"difference" json:"difference"`
	DiscountPercent  decimal.Decimal `db:"discount_percent" json:"discount_percent"`
	Average          decimal.Decimal `db:"average" json:"average"`
	Minimum          int             `db:"minimum" json:"minimum"`
	Maximum          int             `db:"maximum" json:"maximum"`
	MetricEntryCount int             `db:"entries" json:"entries"`
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

func (r *Repository) FindByProductId(productId int64, tx *sql.Tx) (*ProductStats, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select(productStatsTableColumns...).
		From(ProductStatsTableName).
		Where(squirrel.Eq{"product_id": productId})

	var stats ProductStats
	err := q.QueryRow().Scan(
		&stats.ProductId,
		&stats.Price,
		&stats.Difference,
		&stats.DiscountPercent,
		&stats.Average,
		&stats.Minimum,
		&stats.Maximum,
		&stats.MetricEntryCount,
	)

	return &stats, err
}

func (r *Repository) HasProductStats(productId int64, tx *sql.Tx) (bool, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select("COUNT(1)").
		From(ProductStatsTableName).
		Where(squirrel.Eq{"product_id": productId})

	var result int
	err := q.QueryRow().Scan(&result)
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func (r *Repository) CreateProductStats(
	productId int64, price, minimum, maximum, count int, difference, discountPercent, average decimal.Decimal, tx *sql.Tx,
) error {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Insert(ProductStatsTableName).
		Columns(productStatsTableColumns...).
		Values(productId, price, difference, discountPercent, average, minimum, maximum, count)

	_, err := q.Exec()

	return err
}

func (r *Repository) UpdateProductStats(
	productId int64, price, minimum, maximum, count int, difference, discountPercent, average decimal.Decimal, tx *sql.Tx,
) error {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Update(ProductStatsTableName).
		SetMap(map[string]any{
			"price":            price,
			"difference":       difference,
			"discount_percent": discountPercent,
			"average":          average,
			"minimum":          minimum,
			"maximum":          maximum,
			"entries":          count,
		}).
		Where(squirrel.Eq{"product_id": productId})

	_, err := q.Exec()

	return err
}
