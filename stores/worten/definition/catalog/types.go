package catalog

import "github.com/shopspring/decimal"

type Store struct {
	Slug    string
	Name    string
	Website string
}

type StoreProduct struct {
	StoreSlug string
	EAN       []string
	SKU       []string
	Name      string
	Brand     string
	Price     decimal.Decimal
	Available bool
	ImageLink string
	Link      string
}
