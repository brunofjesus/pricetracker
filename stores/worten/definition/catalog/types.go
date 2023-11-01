package catalog

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
	Price     int
	Available bool
	ImageLink string
	Link      string
}
