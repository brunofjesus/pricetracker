package model

const ProductEanTableName = "product_ean"

type ProductEan struct {
	ProductId int64 `db:"product_id"`
	Ean       int64 `db:"ean"`
}
