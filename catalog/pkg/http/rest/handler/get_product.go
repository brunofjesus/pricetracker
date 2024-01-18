package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/brunofjesus/pricetracker/catalog/pkg/http/rest/util"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"github.com/go-chi/chi/v5"
)

func GetProduct(finder *product.Finder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)

		if err != nil {
			util.ErrorJSON(w, errors.New("product does not exist"), http.StatusNotFound)
			return
		}

		productWithMetrics, err := finder.FindProductById(productId)

		if err != nil {
			util.ErrorJSON(w, fmt.Errorf("cannot fetch product: %d", productId), http.StatusInternalServerError)
			return
		}

		util.WriteJSON(w, http.StatusOK, &productWithMetrics)
	}
}
