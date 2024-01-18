package filters

import (
	"fmt"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/common"
	"github.com/brunofjesus/pricetracker/catalog/pkg/pagination"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"net/http"
)

func ToQueryParameters(page pagination.PaginatedQuery, filters product.FinderFilters) string {
	var result = "?"

	// pagination
	result += fmt.Sprintf("page=%d&", page.Page)
	result += fmt.Sprintf("size=%d&", page.PageSize)
	result += fmt.Sprintf("sortField=%s&", page.SortField)
	result += fmt.Sprintf("sortDirection=%s&", page.SortDirection)

	// filters
	for _, id := range filters.ProductId {
		result += fmt.Sprintf("productId=%d&", id)
	}

	if filters.StoreId != nil {
		result += fmt.Sprintf("storeId=%d&", *filters.StoreId)
	}

	if filters.MinPrice != nil {
		result += fmt.Sprintf("minPrice=%f&", *filters.MinPrice)
	}

	if filters.MaxPrice != nil {
		result += fmt.Sprintf("maxPrice=%f&", *filters.MaxPrice)
	}

	if len(filters.NameLike) > 0 {
		result += fmt.Sprintf("name=%s&", filters.NameLike)
	}

	if len(filters.BrandLike) > 0 {
		result += fmt.Sprintf("brand=%s&", filters.BrandLike)
	}

	if filters.Available != nil {
		bAvailable := filters.Available
		val := 0
		if *bAvailable {
			val = 1
		}
		result += fmt.Sprintf("available=%d&", val)
	}

	if len(filters.ProductUrl) > 0 {
		result += fmt.Sprintf("productUrl=%s&", filters.ProductUrl)
	}

	if filters.MinDifference != nil {
		result += fmt.Sprintf("minDifference=%f&", *filters.MinDifference)
	}

	if filters.MaxDifference != nil {
		result += fmt.Sprintf("maxDifference=%f&", *filters.MaxDifference)
	}

	if filters.MinDiscountPercent != nil {
		result += fmt.Sprintf("minDiscountPercent=%f&", *filters.MinDiscountPercent)
	}

	if filters.MaxDiscountPercent != nil {
		result += fmt.Sprintf("maxDiscountPercent=%f&", *filters.MaxDiscountPercent)
	}

	if filters.MinAveragePrice != nil {
		result += fmt.Sprintf("minAveragePrice=%f&", *filters.MinAveragePrice)
	}

	if filters.MaxAveragePrice != nil {
		result += fmt.Sprintf("maxAveragePrice=%f&", *filters.MaxAveragePrice)
	}

	if filters.MinMinimumPrice != nil {
		result += fmt.Sprintf("minMinimumPrice=%f&", *filters.MinMinimumPrice)
	}

	if filters.MaxMinimumPrice != nil {
		result += fmt.Sprintf("minMaximumPrice=%f&", *filters.MaxMinimumPrice)
	}

	if filters.MinMaximumPrice != nil {
		result += fmt.Sprintf("minMaximumPrice=%f&", *filters.MinMaximumPrice)
	}

	if filters.MaxMaximumPrice != nil {
		result += fmt.Sprintf("maxMaximumPrice=%f&", *filters.MaxMaximumPrice)
	}

	return result[:len(result)-1]
}

func FromHttpRequest(r *http.Request) (*product.FinderFilters, error) {
	productId, err := common.GetQueryParamInt64Slice(r, "productId")
	if err != nil {
		return nil, fmt.Errorf("invalid product id: %w", err)
	}

	storeId, err := common.GetQueryParamInt(r, "storeId")
	if err != nil {
		return nil, fmt.Errorf("invalid store id: %w", err)
	}

	minPrice, err := common.GetQueryParamFloat64(r, "minPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min price: %w", err)
	}

	maxPrice, err := common.GetQueryParamFloat64(r, "maxPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid max price: %w", err)
	}

	nameLike := common.GetQueryParam(r, "name")
	brandLike := common.GetQueryParam(r, "brand")
	available := common.GetQueryParamNullBoolean(r, "available")
	productUrl := common.GetQueryParam(r, "productUrl")

	minDifference, err := common.GetQueryParamFloat64(r, "minDifference")
	if err != nil {
		return nil, fmt.Errorf("invalid min difference: %w", err)
	}

	maxDifference, err := common.GetQueryParamFloat64(r, "maxDifference")
	if err != nil {
		return nil, fmt.Errorf("invalid max difference: %w", err)
	}

	minDiscountPercent, err := common.GetQueryParamFloat64(r, "minDiscountPercent")
	if err != nil {
		return nil, fmt.Errorf("invalid min discount: %w", err)
	}

	maxDiscountPercent, err := common.GetQueryParamFloat64(r, "maxDiscountPercent")
	if err != nil {
		return nil, fmt.Errorf("invalid max discount: %w", err)
	}

	minAvgPrice, err := common.GetQueryParamFloat64(r, "minAveragePrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min average price: %w", err)
	}

	maxAvgPrice, err := common.GetQueryParamFloat64(r, "maxAveragePrice")
	if err != nil {
		return nil, fmt.Errorf("invalid max average price: %w", err)
	}

	minMinPrice, err := common.GetQueryParamFloat64(r, "minMinimumPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min minimum price: %w", err)
	}

	maxMinPrice, err := common.GetQueryParamFloat64(r, "maxMinimumPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min maximum price: %w", err)
	}

	minMaxPrice, err := common.GetQueryParamFloat64(r, "minMaximumPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min maximum price: %w", err)
	}

	maxMaxPrice, err := common.GetQueryParamFloat64(r, "maxMaximumPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid max maximum price: %w", err)
	}

	return &product.FinderFilters{
		ProductId:  productId,
		StoreId:    storeId,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		NameLike:   nameLike,
		BrandLike:  brandLike,
		Available:  available,
		ProductUrl: productUrl,

		MinDifference:      minDifference,
		MaxDifference:      maxDifference,
		MinDiscountPercent: minDiscountPercent,
		MaxDiscountPercent: maxDiscountPercent,
		MinAveragePrice:    minAvgPrice,
		MaxAveragePrice:    maxAvgPrice,
		MinMinimumPrice:    minMinPrice,
		MinMaximumPrice:    minMaxPrice,
		MaxMinimumPrice:    maxMinPrice,
		MaxMaximumPrice:    maxMaxPrice,
	}, nil
}
