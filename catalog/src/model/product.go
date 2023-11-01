package model

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
