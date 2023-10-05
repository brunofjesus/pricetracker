package store

import "time"

type WortenSolveURLResponse struct {
	Data struct {
		SolveURL struct {
			Typename string `json:"__typename"`
			ID       string `json:"_id"`
			Code     string `json:"code"`
			Success  bool   `json:"success"`
			Message  any    `json:"message"`
			Item     struct {
				Typename  string `json:"__typename"`
				ID        string `json:"id"`
				Canonical string `json:"canonical"`
				Noindex   bool   `json:"noindex"`
				Type      string `json:"type"`
				Data      struct {
					Config struct {
						Children struct {
							Sort []string `json:"sort"`
						} `json:"children"`
						Attributes []any `json:"attributes"`
						Facebook   struct {
							Description any `json:"description"`
							Image       any `json:"image"`
							Title       any `json:"title"`
						} `json:"facebook"`
						Image []struct {
							ID            int    `json:"id"`
							ProxyPathName string `json:"proxyPathName"`
						} `json:"image"`
						Menu struct {
							Name any `json:"name"`
						} `json:"menu"`
						MetaDescription string `json:"meta_description"`
						PageTitle       string `json:"page_title"`
						Title           string `json:"title"`
						Twitter         struct {
							Description any `json:"description"`
							Image       any `json:"image"`
							Title       any `json:"title"`
						} `json:"twitter"`
					} `json:"config"`
					Entity struct {
						Type     string `json:"type"`
						EntityID string `json:"entity_id"`
						ID       string `json:"id"`
					} `json:"entity"`
					Meta struct {
						Markets   []string `json:"markets"`
						PageType  string   `json:"pageType"`
						Tags      []string `json:"tags"`
						Analytics struct {
							Type string `json:"type"`
						} `json:"analytics"`
						DefaultCategory struct {
							Priority any `json:"priority"`
						} `json:"default_category"`
						IsCategory bool `json:"is_category"`
						Sitemap    struct {
							Frequency any `json:"frequency"`
							Priority  any `json:"priority"`
						} `json:"sitemap"`
					} `json:"meta"`
					Refs []any `json:"refs"`
				} `json:"data"`
				Redirect any `json:"redirect"`
			} `json:"item"`
		} `json:"solveURL"`
	} `json:"data"`
}

type WortenBrowseProductResponse struct {
	Data struct {
		BrowseProducts struct {
			Hits         []WortenProductHit `json:"hits"`
			HasNextPage  bool               `json:"hasNextPage"`
			TotalHits    int                `json:"totalHits"`
			FilterErrors []any              `json:"filterErrors"`
			Typename     string             `json:"__typename"`
		} `json:"browseProducts"`
	} `json:"data"`
}

type WortenProductHit struct {
	ID      string `json:"_id"`
	Product struct {
		ID          string   `json:"id"`
		Ean         []string `json:"ean"`
		Sku         string   `json:"sku"`
		Tags        []string `json:"tags"`
		Name        string   `json:"name"`
		BrandName   string   `json:"brandName"`
		Description string   `json:"description"`
		Type        string   `json:"type"`
		URL         string   `json:"url"`
		Image       struct {
			Type string `json:"type"`
			Data struct {
				Transforms struct {
					Thumbnail string `json:"thumbnail"`
					Default   string `json:"default"`
					Zoom      string `json:"zoom"`
				} `json:"transforms"`
			} `json:"data"`
			URL      string `json:"url"`
			MimeType string `json:"mimeType"`
			Typename string `json:"__typename"`
		} `json:"image"`
		DefaultCategory struct {
			URL      string `json:"url"`
			Typename string `json:"__typename"`
		} `json:"defaultCategory"`
		Typename string `json:"__typename"`
	} `json:"product"`
	TotalOffers int `json:"totalOffers"`
	Data        struct {
		PriceRange      string    `json:"price_range"`
		IsInStock       bool      `json:"is_in_stock"`
		ClassType       string    `json:"class-type"`
		ClassBrickID    string    `json:"class-brick_id"`
		ModelID         string    `json:"model_id"`
		EntityID        string    `json:"entity_id"`
		LowerPriceRange string    `json:"lower_price_range"`
		GroupID         string    `json:"group_id"`
		CategoryPath    []string  `json:"category_path"`
		PromoEnd        time.Time `json:"promo_end"`
		PromoStart      time.Time `json:"promo_start"`
	} `json:"data,omitempty"`
	Campaigns    []any `json:"campaigns"`
	WinningOffer struct {
		ID         string `json:"_id"`
		ChannelID  string `json:"channelId"`
		OfferID    string `json:"offerId"`
		SellerType string `json:"sellerType"`
		Seller     struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Typename string `json:"__typename"`
		} `json:"seller"`
		Price struct {
			Currency string `json:"currency"`
			Value    string `json:"value"`
			Typename string `json:"__typename"`
		} `json:"price"`
		IsInStock bool     `json:"isInStock"`
		Tags      []string `json:"tags"`
		Typename  string   `json:"__typename"`
	} `json:"winningOffer"`
	SecondOfferPrice struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
		Typename string `json:"__typename"`
	} `json:"secondOfferPrice"`
	Typename string `json:"__typename"`
}
