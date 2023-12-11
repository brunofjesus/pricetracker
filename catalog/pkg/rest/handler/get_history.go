package handler

import (
	"fmt"
	"github.com/brunofjesus/pricetracker/catalog/internal/repository/price"
	"github.com/brunofjesus/pricetracker/catalog/pkg/rest/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

type HistoryResponseItem struct {
	DateTime time.Time `json:"date_time"`
	Price    int       `json:"price"`
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)

	from, err := getTimestampFromQueryParam(
		r, "from",
		time.Now().AddDate(0, 0, -30),
	)

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	to, err := getTimestampFromQueryParam(
		r, "to", time.Now(),
	)

	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	prices, err := price.GetPriceRepository().FindPricesBetween(productId, from, to, nil)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	result := make([]HistoryResponseItem, len(prices))

	for i, productPrice := range prices {
		result[i] = HistoryResponseItem{
			DateTime: productPrice.DateTime,
			Price:    productPrice.Price,
		}
	}

	utils.WriteJSON(w, http.StatusOK, result)
}

func getTimestampFromQueryParam(r *http.Request, key string, fallback time.Time) (time.Time, error) {
	seconds, err := utils.GetQueryParamInt64(r, key)

	if err != nil {
		return fallback, fmt.Errorf(
			"invalid timestamp `%s` value `%s`: %w",
			key, utils.GetQueryParam(r, key), err,
		)
	}

	if seconds == -1 {
		return fallback, nil
	}

	return time.Unix(seconds, 0), nil
}
