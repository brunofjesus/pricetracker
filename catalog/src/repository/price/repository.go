package price

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/src/model"
	"github.com/brunofjesus/pricetracker/catalog/src/repository"
	"github.com/shopspring/decimal"
)

type PriceRepository interface {
	GetPrices(productId int64, offset int64, limit int, orderBy, direction string, tx *sql.Tx) ([]model.ProductPrice, error)
	CountPrices(productId int64, tx *sql.Tx) (int64, error)
	CreatePrice(productId int64, price decimal.Decimal, timestamp time.Time, tx *sql.Tx) error
}

type priceRepository struct {
	db *sql.DB
	qb *squirrel.StatementBuilderType
}

func NewPriceRepository(db *sql.DB) PriceRepository {
	return &priceRepository{
		db: db,
		qb: repository.QueryBuilder(db),
	}
}

// GetPrices implements PriceRepository.
func (r *priceRepository) GetPrices(productId int64, offset int64, limit int, orderBy, direction string, tx *sql.Tx) ([]model.ProductPrice, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select("product_id", "date_time", "price").
		From(model.ProductPriceTableName).
		Where(squirrel.Eq{"product_id": productId}).
		OrderBy(fmt.Sprintf("%s %s", orderBy, direction)).
		Offset(uint64(offset)).Limit(uint64(limit))

	var prices []model.ProductPrice
	rows, err := q.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var price model.ProductPrice
		err := rows.Scan(
			&price.ProductId,
			&price.DateTime,
			&price.Price,
		)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}

	return prices, nil
}

// CountPrices implements PriceRepository.
func (r *priceRepository) CountPrices(productId int64, tx *sql.Tx) (int64, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)
	q := qb.Select("COUNT(*)").
		From(model.ProductPriceTableName).
		Where(squirrel.Eq{"product_id": productId})

	var count int64
	err := q.QueryRow().Scan(&count)
	return count, err
}

// CreatePrice implements PriceRepository.
func (r *priceRepository) CreatePrice(productId int64, price decimal.Decimal, timestamp time.Time, tx *sql.Tx) error {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Insert(model.ProductPriceTableName).
		Columns("product_id", "date_time", "price").
		Values(productId, timestamp, price)

	_, err := q.Exec()

	return err
}
