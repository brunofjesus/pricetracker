package util

import (
	"fmt"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/rest/utils"
	"github.com/brunofjesus/pricetracker/catalog/pkg/pagination"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"github.com/brunofjesus/pricetracker/catalog/util/nulltype"
	"net/http"
)

func QueryParamString(page pagination.PaginatedQuery, filters product.FinderFilters) string {
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

	if filters.StoreId > -1 {
		result += fmt.Sprintf("storeId=%d&", filters.StoreId)
	}

	if filters.MinPrice > -1 {
		result += fmt.Sprintf("minPrice=%f&", filters.MinPrice)
	}

	if filters.MaxPrice > -1 {
		result += fmt.Sprintf("maxPrice=%f&", filters.MaxPrice)
	}

	if len(filters.NameLike) > 0 {
		result += fmt.Sprintf("name=%s&", filters.NameLike)
	}

	if len(filters.BrandLike) > 0 {
		result += fmt.Sprintf("brand=%s&", filters.BrandLike)
	}

	if nulltype.IsUndefined(filters.Available) == false {
		result += fmt.Sprintf("available=%d&", filters.Available)
	}

	if len(filters.ProductUrl) > 0 {
		result += fmt.Sprintf("productUrl=%s&", filters.ProductUrl)
	}

	if filters.MinDifference > -1 {
		result += fmt.Sprintf("minDifference=%f&", filters.MinDifference)
	}

	if filters.MaxDifference > -1 {
		result += fmt.Sprintf("maxDifference=%f&", filters.MaxDifference)
	}

	if filters.MinDiscountPercent > -1 {
		result += fmt.Sprintf("minDiscountPercent=%f&", filters.MinDiscountPercent)
	}

	if filters.MaxDiscountPercent > -1 {
		result += fmt.Sprintf("maxDiscountPercent=%f&", filters.MaxDiscountPercent)
	}

	if filters.MinAveragePrice > -1 {
		result += fmt.Sprintf("minAveragePrice=%f&", filters.MinAveragePrice)
	}

	if filters.MaxAveragePrice > -1 {
		result += fmt.Sprintf("maxAveragePrice=%f&", filters.MaxAveragePrice)
	}

	if filters.MinMinimumPrice > -1 {
		result += fmt.Sprintf("minMinimumPrice=%f&", filters.MinMinimumPrice)
	}

	if filters.MaxMinimumPrice > -1 {
		result += fmt.Sprintf("minMaximumPrice=%f&", filters.MinMaximumPrice)
	}

	if filters.MinMaximumPrice > -1 {
		result += fmt.Sprintf("minMaximumPrice=%f&", filters.MaxMaximumPrice)
	}

	if filters.MaxMaximumPrice > -1 {
		result += fmt.Sprintf("maxMaximumPrice=%f&", filters.MaxMaximumPrice)
	}

	return result[:len(result)-1]
}

func GetProductSearchFilters(r *http.Request) (*product.FinderFilters, error) {
	productId, err := utils.GetQueryParamInt64Slice(r, "productId")
	if err != nil {
		return nil, fmt.Errorf("invalid product id: %w", err)
	}

	storeId, err := utils.GetQueryParamInt(r, "storeId")
	if err != nil {
		return nil, fmt.Errorf("invalid store id: %w", err)
	}

	minPrice, err := utils.GetQueryParamFloat64(r, "minPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min price: %w", err)
	}

	maxPrice, err := utils.GetQueryParamFloat64(r, "maxPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid max price: %w", err)
	}

	nameLike := utils.GetQueryParam(r, "name")
	brandLike := utils.GetQueryParam(r, "brand")
	available := utils.GetQueryParamNullBoolean(r, "available")
	productUrl := utils.GetQueryParam(r, "productUrl")

	minDifference, err := utils.GetQueryParamFloat64(r, "minDifference")
	if err != nil {
		return nil, fmt.Errorf("invalid min difference: %w", err)
	}

	maxDifference, err := utils.GetQueryParamFloat64(r, "maxDifference")
	if err != nil {
		return nil, fmt.Errorf("invalid max difference: %w", err)
	}

	minDiscountPercent, err := utils.GetQueryParamFloat64(r, "minDiscountPercent")
	if err != nil {
		return nil, fmt.Errorf("invalid min discount: %w", err)
	}

	maxDiscountPercent, err := utils.GetQueryParamFloat64(r, "maxDiscountPercent")
	if err != nil {
		return nil, fmt.Errorf("invalid max discount: %w", err)
	}

	minAvgPrice, err := utils.GetQueryParamFloat64(r, "minAveragePrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min average price: %w", err)
	}

	maxAvgPrice, err := utils.GetQueryParamFloat64(r, "maxAveragePrice")
	if err != nil {
		return nil, fmt.Errorf("invalid max average price: %w", err)
	}

	minMinPrice, err := utils.GetQueryParamFloat64(r, "minMinimumPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min minimum price: %w", err)
	}

	maxMinPrice, err := utils.GetQueryParamFloat64(r, "maxMinimumPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min maximum price: %w", err)
	}

	minMaxPrice, err := utils.GetQueryParamFloat64(r, "minMaximumPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min maximum price: %w", err)
	}

	maxMaxPrice, err := utils.GetQueryParamFloat64(r, "maxMaximumPrice")
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

func Float64ToString(n float64) string {
	return fmt.Sprintf("%.2f", n)
}
