package model

type MqStore struct {
	Slug    string
	Name    string
	Website string
}

type MqStoreProduct struct {
	StoreSlug string
	EAN       []string
	SKU       []string
	Name      string
	Brand     string
	Price     int
	Available bool
	ImageLink string
	Link      string
}
