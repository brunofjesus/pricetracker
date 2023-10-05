package store

import "time"

type WortenCategoriesResponse struct {
	Data struct {
		SolveURL struct {
			Typename  string `json:"__typename"`
			ID        string `json:"_id"`
			Code      string `json:"code"`
			Success   bool   `json:"success"`
			Message   any    `json:"message"`
			ContextID string `json:"contextId"`
			Layout    struct {
				Typename string `json:"__typename"`
				ID       string `json:"_id"`
				Modules  []struct {
					Typename   string `json:"__typename"`
					Order      int    `json:"order"`
					Priority   int    `json:"priority"`
					TargetedBy string `json:"targetedBy"`
					Targets    string `json:"targets"`
					Refs       []struct {
						Typename string `json:"__typename"`
						Key      string `json:"_key"`
						Type     string `json:"_type"`
						Valid    bool   `json:"valid"`
						ID       string `json:"id"`
						URL      string `json:"url"`
					} `json:"refs"`
				} `json:"modules"`
			} `json:"layout"`
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
