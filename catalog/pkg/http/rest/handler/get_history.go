package handler

import (
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/rest/utils"
	"github.com/brunofjesus/pricetracker/catalog/pkg/price"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

func GetHistory(finder *price.Finder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)

		from, err := utils.GetTimestampFromQueryParam(
			r, "from",
			time.Now().AddDate(0, 0, -30),
		)

		if err != nil {
			utils.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		to, err := utils.GetTimestampFromQueryParam(
			r, "to", time.Now(),
		)

		if err != nil {
			utils.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		result, err := finder.FindPriceHistoryBetween(productId, from, to, nil)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, result)
	}
}
