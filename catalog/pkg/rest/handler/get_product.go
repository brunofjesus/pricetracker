package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"github.com/brunofjesus/pricetracker/catalog/pkg/rest/utils"
	"github.com/go-chi/chi/v5"
)

func GetProduct(finder *product.Finder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)

		if err != nil {
			utils.ErrorJSON(w, errors.New("product does not exist"), http.StatusNotFound)
			return
		}

		productWithMetrics, err := finder.FindProductById(productId)

		if err != nil {
			utils.ErrorJSON(w, fmt.Errorf("cannot fetch product: %d", productId), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, &productWithMetrics)
	}
}
