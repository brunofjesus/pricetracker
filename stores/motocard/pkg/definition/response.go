package definition

type MotocardPageResponse struct {
	Context struct {
		Locale struct {
			Country struct {
				Name                  string `json:"name"`
				Code                  string `json:"code"`
				RealCode              string `json:"realCode"`
				IDRate                int    `json:"idRate"`
				Vat                   int    `json:"vat"`
				VatName               any    `json:"vatName"`
				CustomsIncludedFree   bool   `json:"customsIncludedFree"`
				RequiresGdpr          bool   `json:"requiresGdpr"`
				DefaultLanguageAutoID any    `json:"defaultLanguageAutoId"`
				LanguageAutos         []any  `json:"languageAutos"`
				Hidden                bool   `json:"hidden"`
			} `json:"country"`
			Currency struct {
				Name           string `json:"name"`
				Code           string `json:"code"`
				Symbol         string `json:"symbol"`
				ConversionRate int    `json:"conversionRate"`
				IsVirtual      bool   `json:"isVirtual"`
			} `json:"currency"`
			Language struct {
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"language"`
		} `json:"locale"`
		Popups struct {
			Cookies    bool `json:"cookies"`
			Newsletter struct {
				Show  bool `json:"show"`
				Promo bool `json:"promo"`
			} `json:"newsletter"`
			ChangeCountry      bool `json:"changeCountry"`
			AutoTranslateModal struct {
				Show  bool `json:"show"`
				Texts any  `json:"texts"`
			} `json:"autoTranslateModal"`
			RegisterClubModal struct {
				Show bool `json:"show"`
			} `json:"registerClubModal"`
		} `json:"popups"`
		Auth []any `json:"auth"`
		Cart struct {
			Quantity      int `json:"quantity"`
			Total         int `json:"total"`
			ItemsTotal    int `json:"itemsTotal"`
			ShippingTotal int `json:"shippingTotal"`
			PendingOrders any `json:"pendingOrders"`
		} `json:"cart"`
		Department struct {
			ID                  int    `json:"id"`
			Name                string `json:"name"`
			Slug                string `json:"slug"`
			Identifier          string `json:"identifier"`
			Claim               string `json:"claim"`
			Company             string `json:"company"`
			DomainComercialName string `json:"domainComercialName"`
		} `json:"department"`
		TranslatedURL struct {
			En string `json:"en"`
			Es string `json:"es"`
			Fr string `json:"fr"`
			It string `json:"it"`
			De string `json:"de"`
			Pt string `json:"pt"`
		} `json:"translatedUrl"`
		Canonical     bool   `json:"canonical"`
		Title         string `json:"title"`
		Description   string `json:"description"`
		ServerDate    int    `json:"serverDate"`
		CookieConsent struct {
			Required        bool `json:"required"`
			Statistics      bool `json:"statistics"`
			Personalization bool `json:"personalization"`
			Marketing       bool `json:"marketing"`
		} `json:"cookieConsent"`
		B2B              bool `json:"b2b"`
		IsFranceLaw      bool `json:"isFranceLaw"`
		IsPreBlackFriday bool `json:"isPreBlackFriday"`
		IsCyberMonday    bool `json:"isCyberMonday"`
		Onesignal        struct {
			Uuids                any   `json:"uuids"`
			QueryForSubscription bool  `json:"queryForSubscription"`
			UpdateTags           bool  `json:"updateTags"`
			DeleteTags           bool  `json:"deleteTags"`
			Tags                 []any `json:"tags"`
		} `json:"onesignal"`
		FlashSuccessMessages []any  `json:"flashSuccessMessages"`
		FlashErrorMessages   []any  `json:"flashErrorMessages"`
		IsBlogAvailable      bool   `json:"isBlogAvailable"`
		HasVat               bool   `json:"hasVat"`
		Device               string `json:"device"`
		Version              string `json:"version"`
		SeoTags              []any  `json:"seoTags"`
		PreferredStoreID     any    `json:"preferredStoreId"`
	} `json:"context"`
	Common any `json:"common"`
	View   struct {
		Property     []any                    `json:"property"`
		Category     []any                    `json:"category"`
		Section      []any                    `json:"section"`
		Brand        []any                    `json:"brand"`
		ProductGroup []any                    `json:"product_group"`
		Promo        []any                    `json:"promo"`
		IsOutlet     bool                     `json:"isOutlet"`
		PageOffset   int                      `json:"pageOffset"`
		BannersTop   any                      `json:"bannersTop"`
		BannerList   any                      `json:"bannerList"`
		Results      []MotorcardProductResult `json:"results"`
		TotalHits    int                      `json:"totalHits"`
		Paginator    struct {
			CurrentPage           int    `json:"currentPage"`
			Pages                 int    `json:"pages"`
			TotalItems            int    `json:"totalItems"`
			TemplatePageURL       string `json:"templatePageUrl"`
			FirstPageURL          string `json:"firstPageUrl"`
			OtherPagesTemplateURL string `json:"otherPagesTemplateUrl"`
		} `json:"paginator"`
	} `json:"view"`
}

type MotorcardProductResult struct {
	ID          int      `json:"id"`
	PropertyIds []int    `json:"propertyIds"`
	Flags       []string `json:"flags"`
	Name        string   `json:"name"`
	ProductCode int      `json:"productCode"`
	Title       struct {
		Category string `json:"category"`
		Brand    string `json:"brand"`
		Name     string `json:"name"`
		Full     string `json:"full"`
	} `json:"title"`
	ShortDescription struct {
		Auto int    `json:"auto"`
		Text string `json:"text"`
	} `json:"shortDescription"`
	References struct {
		Motocard string `json:"motocard"`
		Google   string `json:"google"`
		Criteo   string `json:"criteo"`
		External string `json:"external"`
	} `json:"references"`
	Brand struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Logo string `json:"logo"`
		URL  []any  `json:"url"`
	} `json:"brand"`
	URL     string `json:"url"`
	Promos  []any  `json:"promos"`
	Reviews struct {
		Count         int `json:"count"`
		AverageRating int `json:"averageRating"`
	} `json:"reviews"`
	IsDiscontinued    bool   `json:"isDiscontinued"`
	IsPreorder        bool   `json:"isPreorder"`
	IsSale            bool   `json:"isSale"`
	IsBlackFriday     bool   `json:"isBlackFriday"`
	NumSizesWithStock int    `json:"numSizesWithStock"`
	HasMotorbikes     bool   `json:"hasMotorbikes"`
	Image             string `json:"image"`
	VolumetricWeight  int    `json:"volumetricWeight"`
	Sizes             []struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Prices []struct {
			ID             int    `json:"id"`
			ProductSizeID  int    `json:"product_size_id"`
			RateID         int    `json:"rate_id"`
			ExternalRateID int    `json:"external_rate_id"`
			Type           string `json:"type"`
			BasePrice      int    `json:"base_price"`
			Discount       int    `json:"discount"`
			Cost           int    `json:"cost"`
			CreatedAt      string `json:"created_at"`
			UpdatedAt      string `json:"updated_at"`
			DeletedAt      any    `json:"deleted_at"`
		} `json:"prices"`
		Stock  int    `json:"stock"`
		Gtin   string `json:"gtin"`
		Stocks []struct {
			Stock       int  `json:"stock"`
			WarehouseID int  `json:"warehouseId"`
			IsStore     bool `json:"isStore"`
		} `json:"stocks"`
		HasProviderStock    bool `json:"hasProviderStock"`
		DirectShippingStock struct {
			Num193 int `json:"193"`
			Num198 int `json:"198"`
			Num204 int `json:"204"`
			Num208 int `json:"208"`
		} `json:"directShippingStock"`
		HasStock   bool `json:"hasStock"`
		References struct {
			Motocard string `json:"motocard"`
		} `json:"references"`
	} `json:"sizes"`
	SizesString    string `json:"sizesString"`
	SizingID       int    `json:"sizingId"`
	CreatedAt      string `json:"createdAt"`
	CartGrouping   bool   `json:"cartGrouping"`
	IsDigital      bool   `json:"isDigital"`
	AmountGiftCard int    `json:"amountGiftCard"`
	SkipVat        bool   `json:"skip_vat"`
	SpecialTag     struct {
		Text             string `json:"text"`
		IsFlashPackOffer bool   `json:"is_flash_pack_offer"`
	} `json:"specialTag"`
	DepartmentsName string `json:"departmentsName"`
	SectionsName    string `json:"sectionsName"`
	CategoriesName  string `json:"categoriesName"`
	TypologiesName  string `json:"typologiesName"`
	CustomizationID any    `json:"customizationId"`
	Images          []struct {
		URL   string `json:"url"`
		Alert []any  `json:"alert"`
	} `json:"images"`
	UsagesName     []string `json:"usagesName"`
	CacheVersion   int      `json:"cacheVersion"`
	BasePrice      string   `json:"basePrice"`
	Price          string   `json:"price"`
	Discount       int      `json:"discount"`
	RawBasePrice   int      `json:"rawBasePrice"`
	RawPrice       int      `json:"rawPrice"`
	DiscountAmount string   `json:"discountAmount"`
	ProductOverlay any      `json:"productOverlay"`
	CountdownDate  any      `json:"countdownDate"`
	Tags           []struct {
		IsDiscountTag   bool   `json:"is_discount_tag"`
		Title           string `json:"title"`
		Type            string `json:"type"`
		StatisticsTitle string `json:"statistics-title"`
	} `json:"tags"`
	BestPromoCode   any    `json:"bestPromoCode"`
	HasFreeShipping bool   `json:"hasFreeShipping"`
	RequestCheck    string `json:"requestCheck"`
	Wished          bool   `json:"wished"`
}
