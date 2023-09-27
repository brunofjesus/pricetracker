package projection

import "github.com/shopspring/decimal"

var ProductWithMetricsViewName = "product_metrics"

type ProductWithMetrics struct {
	ProductId  int64           `db:"product_id"`
	StoreId    int64           `db:"store_id"`
	Name       string          `db:"name"`
	Brand      string          `db:"brand"`
	Price      decimal.Decimal `db:"price"`
	Available  bool            `db:"available"`
	ImageUrl   string          `db:"image_url"`
	ProductUrl string          `db:"product_url"`

	Discount         decimal.Decimal `db:"discount"`
	DiscountPercent  decimal.Decimal `db:"discount_percent"`
	Average          decimal.Decimal `db:"average"`
	Maximum          decimal.Decimal `db:"maximum"`
	Minimum          decimal.Decimal `db:"minimum"`
	MetricEntryCount decimal.Decimal `db:"entries"`
	MetricDataSince  decimal.Decimal `db:"metrics_since"`
}
