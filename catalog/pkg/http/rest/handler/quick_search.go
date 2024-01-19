package handler

import (
	"fmt"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/common"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/rest/util"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"net/http"
)

func QuickSearch(finder *product.Finder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := common.GetQueryParam(r, "q")

		if len(q) == 0 {
			util.ErrorJSON(w, fmt.Errorf("q cannot be empty"), http.StatusBadRequest)
			return
		}

		results, err := finder.QuickSearch(q, true)
		if err != nil {
			util.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		util.WriteJSON(w, http.StatusOK, results)
	}
}
