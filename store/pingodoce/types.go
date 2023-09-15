package pingodoce

type PingoDoceSearchResult struct {
	Sections struct {
		Null struct {
			Total    int `json:"total"`
			Products []struct {
				Index  string `json:"_index"`
				Type   string `json:"_type"`
				ID     string `json:"_id"`
				Score  int    `json:"_score"`
				Source struct {
					FirstName        string  `json:"firstName"`
					SecondName       string  `json:"secondName"`
					ThirdName        string  `json:"thirdName"`
					LongDescription  string  `json:"longDescription"`
					ShortDescription string  `json:"shortDescription"`
					Sku              string  `json:"sku"`
					ImagesNumber     int     `json:"imagesNumber"`
					GrossWeight      float64 `json:"grossWeight"`
					Capacity         string  `json:"capacity"`
					NetContent       int     `json:"netContent"`
					NetContentUnit   string  `json:"netContentUnit"`
					AverageWeight    int     `json:"averageWeight"`
					OnlineStatus     string  `json:"onlineStatus"`
					Status           string  `json:"status"`
					Slug             string  `json:"slug"`
					Tags             []any   `json:"tags"`
					Categories       []struct {
						Name string `json:"name"`
						ID   string `json:"id"`
					} `json:"categories"`
					Eans  []string `json:"eans"`
					Brand struct {
						Name string `json:"name"`
						ID   string `json:"id"`
					} `json:"brand"`
					CatalogueID     string   `json:"catalogueId"`
					CategoriesArray []string `json:"categoriesArray"`
					LeafCategories  []struct {
						Name string `json:"name"`
						ID   string `json:"id"`
					} `json:"leafCategories"`
					IsPerishable             bool     `json:"isPerishable"`
					AncestorsCategoriesArray []string `json:"ancestorsCategoriesArray"`
					RegularPrice             float64  `json:"regularPrice"`
					CampaignPrice            int      `json:"campaignPrice"`
					BuyingPrice              float64  `json:"buyingPrice"`
					UnitPrice                float64  `json:"unitPrice"`
					Promotion                struct {
						BeginDate  any `json:"beginDate"`
						Amount     any `json:"amount"`
						PayAmount  any `json:"payAmount"`
						EndDate    any `json:"endDate"`
						TakeAmount any `json:"takeAmount"`
						Type       any `json:"type"`
					} `json:"promotion"`
					MinimumOrderableQuantity int    `json:"minimumOrderableQuantity"`
					MaximumOrderableQuantity int    `json:"maximumOrderableQuantity"`
					CountriesOfOrigin        []any  `json:"countriesOfOrigin"`
					AdditionalInfo           string `json:"additionalInfo"`
					DurabilityDays           int    `json:"durabilityDays"`
					ActivePromotion          bool   `json:"activePromotion"`
					Advertising              []any  `json:"advertising"`
					SuggestCatalogue         []struct {
						Input    string `json:"input"`
						Weight   int    `json:"weight"`
						Contexts struct {
							Catalogue string `json:"catalogue"`
						} `json:"contexts"`
					} `json:"suggest_catalogue"`
					UnitNoVatPrice     float64 `json:"unitNoVatPrice"`
					CampaignNoVatPrice int     `json:"campaignNoVatPrice"`
					QualitativeIcons   []any   `json:"qualitativeIcons"`
					CapacityType       any     `json:"capacityType"`
					NoVatPrice         float64 `json:"noVatPrice"`
					VatTax             float64 `json:"vatTax"`
					BuyingNoVatPrice   float64 `json:"buyingNoVatPrice"`
					LowestBuyingPrice  float64 `json:"lowestBuyingPrice"`
				} `json:"_source"`
			} `json:"products"`
			Order int `json:"order"`
			Name  any `json:"name"`
		} `json:"null"`
	} `json:"sections"`
	Categories []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"categories"`
	Brands []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"brands"`
}
