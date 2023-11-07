package model

import (
	"time"

	"github.com/shopspring/decimal"
)

const StoreTableName = "store"
const ProductTableName = "product"
const ProductEanTableName = "product_ean"
const ProductSkuTableName = "product_sku"
const ProductPriceTableName = "product_price"
const ProductWithMetricsViewName = "product_metrics"

type Store struct {
	StoreId int64  `db:"store_id"`
	Slug    string `db:"slug"`
	Name    string `db:"name"`
	Website string `db:"website"`
	Active  bool   `db:"active"`
}

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

type ProductEan struct {
	ProductId int64 `db:"product_id"`
	Ean       int64 `db:"ean"`
}

type ProductSku struct {
	ProductId int64  `db:"product_id"`
	Sku       string `db:"sku"`
}

type ProductPrice struct {
	ProductId int64     `db:"product_id"`
	DateTime  time.Time `db:"date_time"`
	Price     int       `db:"price"`
}

type ProductWithMetrics struct {
	ProductId  int64  `db:"product_id"`
	StoreId    int64  `db:"store_id"`
	Name       string `db:"name"`
	Brand      string `db:"brand"`
	Price      int    `db:"price"`
	Available  bool   `db:"available"`
	ImageUrl   string `db:"image_url"`
	ProductUrl string `db:"product_url"`

	Difference       decimal.Decimal `db:"diff"`
	DiscountPercent  decimal.Decimal `db:"discount_percent"`
	Average          decimal.Decimal `db:"average"`
	Maximum          decimal.Decimal `db:"maximum"`
	Minimum          decimal.Decimal `db:"minimum"`
	MetricEntryCount decimal.Decimal `db:"entries"`
	MetricDataSince  decimal.Decimal `db:"metrics_since"`
}
