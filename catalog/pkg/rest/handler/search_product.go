package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"github.com/brunofjesus/pricetracker/catalog/pkg/rest/utils"
	"github.com/brunofjesus/pricetracker/catalog/util/nulltype"
	"github.com/brunofjesus/pricetracker/catalog/util/pagination"
)

func SearchProduct(w http.ResponseWriter, r *http.Request) {

	filters, err := getFilters(r)

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	paginationQuery, err := pagination.FromHttpRequest(r)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	products, err := product.GetMetricsFinder().FindProducts(
		*paginationQuery, *filters,
	)

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

func getFilters(r *http.Request) (*product.ProductFinderFilters, error) {
	productId, err := getQueryParamInt64Slice(r, "productId")
	if err != nil {
		return nil, fmt.Errorf("invalid product id: %w", err)
	}

	storeId, err := getQueryParamInt(r, "storeId")
	if err != nil {
		return nil, fmt.Errorf("invalid store id: %w", err)
	}

	minPrice, err := getQueryParamFloat64(r, "minPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min price: %w", err)
	}

	maxPrice, err := getQueryParamFloat64(r, "maxPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid max price: %w", err)
	}

	nameLike := getQueryParam(r, "name")
	brandLike := getQueryParam(r, "brand")
	available := getQueryParamNullBoolean(r, "available")
	productUrl := getQueryParam(r, "productUrl")

	minDifference, err := getQueryParamFloat64(r, "minDifference")
	if err != nil {
		return nil, fmt.Errorf("invalid min difference: %w", err)
	}

	maxDifference, err := getQueryParamFloat64(r, "maxDifference")
	if err != nil {
		return nil, fmt.Errorf("invalid max difference: %w", err)
	}

	minDiscountPercent, err := getQueryParamFloat64(r, "minDiscountPercent")
	if err != nil {
		return nil, fmt.Errorf("invalid min discount: %w", err)
	}

	maxDiscountPercent, err := getQueryParamFloat64(r, "maxDiscountPercent")
	if err != nil {
		return nil, fmt.Errorf("invalid max discount: %w", err)
	}

	minAvgPrice, err := getQueryParamFloat64(r, "minAveragePrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min average price: %w", err)
	}

	maxAvgPrice, err := getQueryParamFloat64(r, "maxAveragePrice")
	if err != nil {
		return nil, fmt.Errorf("invalid max average price: %w", err)
	}

	minMinPrice, err := getQueryParamFloat64(r, "minMinimumPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min minimum price: %w", err)
	}

	maxMinPrice, err := getQueryParamFloat64(r, "maxMinimumPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min maximum price: %w", err)
	}

	minMaxPrice, err := getQueryParamFloat64(r, "minMaximumPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid min maximum price: %w", err)
	}

	maxMaxPrice, err := getQueryParamFloat64(r, "maxMaximumPrice")
	if err != nil {
		return nil, fmt.Errorf("invalid max maximum price: %w", err)
	}

	return &product.ProductFinderFilters{
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

func getQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func getQueryParamNullBoolean(r *http.Request, key string) nulltype.NullBoolean {
	intQueryParam, err := getQueryParamInt(r, key)

	if err != nil {
		return nulltype.UndefinedValue
	}

	return nulltype.FromInt(intQueryParam)
}

func getQueryParamInt(r *http.Request, key string) (int, error) {
	strVal := getQueryParam(r, key)

	if len(strVal) == 0 {
		return -1, nil
	}

	return strconv.Atoi(strVal)
}

func getQueryParamInt64Slice(r *http.Request, key string) ([]int64, error) {
	queryParamValue := getQueryParam(r, key)

	if !strings.Contains(queryParamValue, ",") {
		intVal, err := getQueryParamInt64(r, key)
		if err != nil || intVal == -1 {
			return []int64{}, err
		}
		return []int64{intVal}, nil
	}

	queryParamValueSlice := strings.Split(queryParamValue, ",")
	result := make([]int64, 0, len(queryParamValueSlice))

	for _, str := range queryParamValueSlice {
		intVal, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return []int64{}, err
		}
		result = append(result, intVal)
	}

	return result, nil
}

func getQueryParamInt64(r *http.Request, key string) (int64, error) {
	strVal := getQueryParam(r, key)

	if len(strVal) == 0 {
		return -1, nil
	}

	return strconv.ParseInt(strVal, 10, 64)
}

func getQueryParamFloat64(r *http.Request, key string) (float64, error) {
	strVal := getQueryParam(r, key)

	if len(strVal) == 0 {
		return -1, nil
	}

	return strconv.ParseFloat(strVal, 64)
}
