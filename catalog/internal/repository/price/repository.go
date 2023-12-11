package price

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/internal/repository"
)

var once sync.Once
var instance PriceRepository

const ProductPriceTableName = "product_price"

type ProductPrice struct {
	ProductId int64     `db:"product_id"`
	DateTime  time.Time `db:"date_time"`
	Price     int       `db:"price"`
}

type PriceRepository interface {
	FindLatestPrice(productId int64, tx *sql.Tx) (*ProductPrice, error)
	FindPricesBetween(productId int64, from time.Time, to time.Time, tx *sql.Tx) ([]ProductPrice, error)
	FindPrices(productId int64, offset int64, limit int, orderBy, direction string, tx *sql.Tx) ([]ProductPrice, error)
	CountPrices(productId int64, tx *sql.Tx) (int64, error)
	CreatePrice(productId int64, price int, timestamp time.Time, tx *sql.Tx) error
}

type priceRepository struct {
	db *sql.DB
	qb *squirrel.StatementBuilderType
}

func GetPriceRepository() PriceRepository {
	once.Do(func() {
		db := repository.GetDatabaseConnection()

		instance = &priceRepository{
			db: db,
			qb: repository.QueryBuilder(db),
		}
	})

	return instance
}

// FindLatestPrice implements PriceRepository.
func (r *priceRepository) FindLatestPrice(productId int64, tx *sql.Tx) (*ProductPrice, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select("product_id", "date_time", "price").
		From(ProductPriceTableName).
		Where(squirrel.Eq{"product_id": productId}).
		OrderBy("date_time desc").
		Offset(0).Limit(1)

	var productPrice ProductPrice
	err := q.QueryRow().Scan(
		&productPrice.ProductId,
		&productPrice.DateTime,
		&productPrice.Price,
	)

	return &productPrice, err
}

// FindPricesBetween implements PriceRepository
func (r *priceRepository) FindPricesBetween(productId int64, from time.Time, to time.Time, tx *sql.Tx) ([]ProductPrice, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select("product_id", "date_time", "price").
		From(ProductPriceTableName).
		Where(
			squirrel.And{
				squirrel.Eq{
					"product_id": productId,
				},
				squirrel.GtOrEq{
					"date_time": from.Format(time.RFC3339),
				},
				squirrel.LtOrEq{
					"date_time": to.Format(time.RFC3339),
				},
			},
		)

	var prices []ProductPrice
	rows, err := q.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var price ProductPrice
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

// FindPrices implements PriceRepository.
func (r *priceRepository) FindPrices(productId int64, offset int64, limit int, orderBy, direction string, tx *sql.Tx) ([]ProductPrice, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select("product_id", "date_time", "price").
		From(ProductPriceTableName).
		Where(squirrel.Eq{"product_id": productId}).
		OrderBy(fmt.Sprintf("%s %s", orderBy, direction)).
		Offset(uint64(offset)).Limit(uint64(limit))

	var prices []ProductPrice
	rows, err := q.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var price ProductPrice
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
		From(ProductPriceTableName).
		Where(squirrel.Eq{"product_id": productId})

	var count int64
	err := q.QueryRow().Scan(&count)
	return count, err
}

// CreatePrice implements PriceRepository.
func (r *priceRepository) CreatePrice(productId int64, price int, timestamp time.Time, tx *sql.Tx) error {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Insert(ProductPriceTableName).
		Columns("product_id", "date_time", "price").
		Values(productId, timestamp, price)

	_, err := q.Exec()

	return err
}
