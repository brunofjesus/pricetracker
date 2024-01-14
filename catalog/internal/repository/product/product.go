package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/internal/repository"
)

const ProductTableName = "product"

type Product struct {
	ProductId  int64  `db:"product_id"`
	StoreId    int64  `db:"store_id"`
	Name       string `db:"name"`
	Brand      string `db:"brand"`
	Price      int    `db:"price"`
	Available  bool   `db:"available"`
	ImageUrl   string `db:"image_url"`
	ProductUrl string `db:"product_url"`
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

func (r *Repository) FindProductById(id int64, tx *sql.Tx) (*Product, error) {
	return r.findOne(tx, squirrel.Eq{"product_id": id})
}

func (r *Repository) FindProductByUrl(url string, tx *sql.Tx) (*Product, error) {
	return r.findOne(tx, squirrel.Eq{"product_url": url})
}

func (r *Repository) CreateProduct(
	storeId int64, name string, brand string, imageUrl string, productUrl string,
	price int, available bool, tx *sql.Tx,
) (int64, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Insert(ProductTableName).
		Columns("store_id", "name", "brand", "price", "available", "image_url", "product_url").
		Values(storeId, name, brand, price, available, imageUrl, productUrl).
		Suffix("RETURNING product_id")

	var id int64
	err := q.QueryRow().Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r *Repository) UpdateProduct(
	productId int64, name string, brand string, imageUrl string, productUrl string,
	price int, available bool, tx *sql.Tx,
) error {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Update(ProductTableName).
		SetMap(map[string]any{
			"name":        name,
			"brand":       brand,
			"image_url":   imageUrl,
			"product_url": productUrl,
			"price":       price,
			"available":   available,
		}).
		Where(squirrel.Eq{"product_id": productId})

	_, err := q.Exec()

	return err
}

func (r *Repository) findOne(tx *sql.Tx, where any, args ...any) (*Product, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Select("product_id", "store_id", "name", "brand", "price", "available", "image_url", "product_url").
		From(ProductTableName).
		Where(where, args...)

	var product Product
	err := q.QueryRow().Scan(
		&product.ProductId,
		&product.StoreId,
		&product.Name,
		&product.Brand,
		&product.Price,
		&product.Available,
		&product.ImageUrl,
		&product.ProductUrl,
	)

	return &product, err
}
