package handler

import (
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/common/filters"
	"net/http"

	"github.com/brunofjesus/pricetracker/catalog/pkg/http/common"

	"github.com/brunofjesus/pricetracker/catalog/pkg/http/rest/util"
	"github.com/brunofjesus/pricetracker/catalog/pkg/pagination"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
)

func SearchProduct(finder *product.Finder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestFilters, err := filters.FromHttpRequest(r)
		if err != nil {
			util.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		fetchEan := common.GetQueryParamBoolean(r, "fetchEan", false)
		fetchSku := common.GetQueryParamBoolean(r, "fetchSku", false)

		paginationQuery, err := pagination.FromHttpRequest(r, 10)
		if err != nil {
			util.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		products, err := finder.FindDetailedProducts(
			*paginationQuery, *requestFilters, fetchEan, fetchSku,
		)

		if err != nil {
			util.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		util.WriteJSON(w, http.StatusOK, products)
	}
}
