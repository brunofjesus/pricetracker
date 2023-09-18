package datastruct

const ProductSkuTableName = "product_sku"

type ProductSku struct {
	ProductId int64  `db:"product_id"`
	Sku       string `db:"sku"`
}
