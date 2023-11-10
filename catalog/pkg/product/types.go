package product

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
