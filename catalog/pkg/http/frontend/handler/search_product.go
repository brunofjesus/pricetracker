package handler

import (
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/frontend/util"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/frontend/view"
	"github.com/brunofjesus/pricetracker/catalog/pkg/pagination"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"net/http"
)

func SearchProduct(finder *product.Finder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters, err := util.GetProductSearchFilters(r)
		if err != nil {
			writeBadRequest(w)
			return
		}

		paginationQuery, err := pagination.FromHttpRequest(r)
		if err != nil {
			writeBadRequest(w)
			return
		}

		if !paginationQuery.HasSortField() {
			paginationQuery.SortField = "discount_percent"
		}

		products, err := finder.FindDetailedProducts(
			*paginationQuery, *filters,
		)

		if err != nil || products == nil {
			writeInternalError(w)
			return
		}

		viewProps := view.ProductsViewProps{
			Page:      *products,
			PageQuery: *paginationQuery,
			Filters:   *filters,
		}

		err = view.ProductsView(viewProps).Render(r.Context(), w)
		if err != nil {
			writeInternalError(w)
			return
		}
	}

}
