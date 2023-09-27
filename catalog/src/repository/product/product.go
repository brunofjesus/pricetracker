package product

import (
	"database/sql"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/brunofjesus/pricetracker/catalog/src/model"
	"github.com/brunofjesus/pricetracker/catalog/src/repository"
	"github.com/shopspring/decimal"
)

var productOnce sync.Once
var productInstance ProductRepository

type ProductRepository interface {
	FindProductById(id int64, tx *sql.Tx) (*model.Product, error)
	FindProductByUrl(url string, tx *sql.Tx) (*model.Product, error)
	CreateProduct(storeId int64, name, brand, imageUrl, productUrl string, price decimal.Decimal, available bool, tx *sql.Tx) (int64, error)
	UpdateProduct(productId int64, name, brand, imageUrl, productUrl string, price decimal.Decimal, available bool, tx *sql.Tx) error
}

type productRepository struct {
	db *sql.DB
	qb *squirrel.StatementBuilderType
}

func GetProductRepository() ProductRepository {
	productOnce.Do(func() {
		db := repository.GetDatabaseConnection()

		productInstance = &productRepository{
			db: db,
			qb: repository.QueryBuilder(db),
		}
	})
	return productInstance
}

// FindProductById implements ProductRepository.
func (r *productRepository) FindProductById(id int64, tx *sql.Tx) (*model.Product, error) {
	return r.findOne(tx, squirrel.Eq{"product_id": id})
}

// FindProductByUrl implements ProductRepository.
func (r *productRepository) FindProductByUrl(url string, tx *sql.Tx) (*model.Product, error) {
	return r.findOne(tx, squirrel.Eq{"product_url": url})
}

// CreateProduct implements ProductRepository.
func (r *productRepository) CreateProduct(
	storeId int64, name string, brand string, imageUrl string, productUrl string,
	price decimal.Decimal, available bool, tx *sql.Tx,
) (int64, error) {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Insert(model.ProductTableName).
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

// UpdateProduct implements ProductRepository.
func (r *productRepository) UpdateProduct(
	productId int64, name string, brand string, imageUrl string, productUrl string,
	price decimal.Decimal, available bool, tx *sql.Tx,
) error {
	qb := repository.QueryBuilderOrDefault(tx, r.qb)

	q := qb.Update(model.ProductTableName).
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

func (r *productRepository) findOne(tx *sql.Tx, where any, args ...any) (*model.Product, error) {
	qb := r.qb
	if tx != nil {
		qb = repository.QueryBuilder(tx)
	}

	q := qb.Select("product_id", "store_id", "name", "brand", "price", "available", "image_url", "product_url").
		From(model.ProductTableName).
		Where(where, args...)

	var product model.Product
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
