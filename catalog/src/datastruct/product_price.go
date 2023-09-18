package datastruct

import (
	"time"

	"github.com/shopspring/decimal"
)

const ProductPriceTableName = "product_price"

type ProductPrice struct {
	ProductId int64           `db:"product_id"`
	DateTime  time.Time       `db:"date_time"`
	Price     decimal.Decimal `db:"price"`
}
