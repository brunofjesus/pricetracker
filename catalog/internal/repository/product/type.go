package product

import "github.com/shopspring/decimal"

type Product struct {
	ProductId  int64       `db:"product_id" json:"product_id"`
	StoreId    int64       `db:"store_id" json:"store_id,omitempty"`
	Name       string      `db:"name" json:"name"`
	Brand      string      `db:"brand" json:"brand"`
	Price      int         `db:"price" json:"price"`
	Currency   string      `db:"currency" json:"currency"`
	Available  bool        `db:"available" json:"available"`
	ImageUrl   string      `db:"image_url" json:"image_url"`
	ProductUrl string      `db:"product_url" json:"product_url"`
	Store      *Store      `json:"store,omitempty"`
	Statistics *Statistics `json:"statistics,omitempty"`
	Ean        []int64     `json:"ean,omitempty"`
	Sku        []string    `json:"sku,omitempty"`
}

type Store struct {
	StoreId int64  `db:"store_id" json:"storeId,omitempty"`
	Name    string `db:"store_name" json:"name"`
	Slug    string `db:"store_slug" json:"slug"`
	Website string `db:"store_website" json:"website"`
}

type Statistics struct {
	ProductId        int64           `db:"product_id" json:"product_id,omitempty"`
	Difference       decimal.Decimal `db:"difference" json:"difference"`
	DiscountPercent  decimal.Decimal `db:"discount_percent" json:"discount_percent"`
	Average          decimal.Decimal `db:"average" json:"average"`
	Minimum          decimal.Decimal `db:"minimum" json:"minimum"`
	Maximum          decimal.Decimal `db:"maximum" json:"maximum"`
	MetricEntryCount decimal.Decimal `db:"entries" json:"entries"`
}

type ProductEan struct {
	ProductId int64 `db:"product_id"`
	Ean       int64 `db:"ean"`
}

type ProductSku struct {
	ProductId int64  `db:"product_id"`
	Sku       string `db:"sku"`
}
