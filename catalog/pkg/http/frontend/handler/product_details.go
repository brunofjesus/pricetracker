package handler

import (
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/frontend/view"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/rest/utils"
	"github.com/brunofjesus/pricetracker/catalog/pkg/price"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

func ProductDetails(productFinder *product.Finder, priceFinder *price.Finder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productId, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)
		if err != nil {
			writeBadRequest(w)
			return
		}

		from, err := utils.GetTimestampFromQueryParam(
			r, "from",
			time.Now().AddDate(0, 0, -30),
		)

		if err != nil {
			writeBadRequest(w)
			return
		}

		to, err := utils.GetTimestampFromQueryParam(
			r, "to", time.Now(),
		)

		if err != nil {
			writeBadRequest(w)
			return
		}

		item, err := productFinder.FindProductById(productId)
		if err != nil || item == nil {
			writeInternalError(w)
			return
		}

		prices, err := priceFinder.FindPriceHistoryBetween(productId, from, to, nil)
		if err != nil {
			writeInternalError(w)
			return
		}

		viewProps := view.DetailsViewProps{
			Product: *item,
			Prices:  prices,
		}

		err = view.DetailsView(viewProps).Render(r.Context(), w)
		if err != nil {
			writeInternalError(w)
			return
		}
	}
}
