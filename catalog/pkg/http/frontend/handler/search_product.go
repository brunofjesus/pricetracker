package handler

import (
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/common/filters"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/frontend/view"
	"github.com/brunofjesus/pricetracker/catalog/pkg/pagination"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"github.com/brunofjesus/pricetracker/catalog/pkg/store"
	"net/http"
)

func SearchProduct(productFinder *product.Finder, storeFinder *store.Finder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestFilters, err := filters.FromHttpRequest(r)
		if err != nil {
			writeBadRequest(w)
			return
		}

		paginationQuery, err := pagination.FromHttpRequest(r, 12)
		if err != nil {
			writeBadRequest(w)
			return
		}

		if !paginationQuery.HasSortField() {
			paginationQuery.SortField = "discount_percent"
		}

		products, err := productFinder.FindDetailedProducts(
			*paginationQuery, *requestFilters, false, false,
		)

		if err != nil || products == nil {
			writeInternalError(w)
			return
		}

		stores, err := storeFinder.FindStores()

		if err != nil {
			writeInternalError(w)
			return
		}

		viewProps := view.ProductsViewProps{
			Page:      *products,
			PageQuery: *paginationQuery,
			Filters:   *requestFilters,
			Stores:    stores,
		}

		err = view.ProductsView(viewProps).Render(r.Context(), w)
		if err != nil {
			writeInternalError(w)
			return
		}
	}

}
