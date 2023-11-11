package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/brunofjesus/pricetracker/catalog/internal/repository/product"
	"github.com/brunofjesus/pricetracker/catalog/pkg/rest/utils"
	"github.com/go-chi/chi/v5"
)

func GetProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)

	if err != nil {
		utils.ErrorJSON(w, errors.New("product does not exist"), http.StatusNotFound)
		return
	}

	metricsRepository := product.GetProductMetricsRepository()

	productWithMetrics, err := metricsRepository.FindProductById(productId, nil)

	if err != nil {
		utils.ErrorJSON(w, fmt.Errorf("cannot fetch product: %d", productId), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, 200, &productWithMetrics)
}
