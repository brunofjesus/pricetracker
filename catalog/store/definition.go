package store

import "github.com/shopspring/decimal"

type Store interface {
	Id() string
	Name() string
	Website() string
	Crawl(productChannel chan StoreProduct)
}

type StoreProduct struct {
	StoreId   string
	StoreName string
	EAN       []string
	SKU       []string
	Name      string
	Brand     string
	Price     decimal.Decimal
	Available bool
	ImageLink string
	Link      string
}
