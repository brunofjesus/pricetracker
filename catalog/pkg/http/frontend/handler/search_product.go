package handler

import (
	"fmt"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/frontend/view"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/rest/utils"
	"github.com/brunofjesus/pricetracker/catalog/pkg/pagination"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"net/http"
)

func SearchProduct(finder *product.Finder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters, err := getProductSearchFilters(r)
		if err != nil {
			writeBadRequest(w)
			return
		}

		paginationQuery, err := pagination.FromHttpRequest(r)
		if err != nil {
			writeBadRequest(w)
			return
		}

		products, err := finder.FindDetailedProducts(
			*paginationQuery, *filters,
		)

		if err != nil || products == nil {
			writeInternalError(w)
			return
		}

		err = view.ProductsView(*products).Render(r.Context(), w)
		if err != nil {
			writeInternalError(w)
			return
		}
	}

}

// TODO: duplicate from rest package
func getProductSearchFilters(r *http.Request) (*product.FinderFilters, error) {
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
