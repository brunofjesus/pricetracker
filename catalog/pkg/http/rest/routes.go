package rest

import (
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/rest/handler"
	"github.com/brunofjesus/pricetracker/catalog/pkg/price"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"github.com/go-chi/chi/v5"
)

type RouteFunc = func(r chi.Router)

type V1ApiProps struct {
	ProductFinder *product.Finder
	PriceFinder   *price.Finder
}

func v1Routes(p V1ApiProps) RouteFunc {
	return func(r chi.Router) {
		r.Get("/product/{productId}", handler.GetProduct(p.ProductFinder))
		r.Get("/product/search", handler.SearchProduct(p.ProductFinder))
		r.Get("/product/{productId}/history", handler.GetHistory(p.PriceFinder))
	}
}

func AddRoutes(r chi.Router, p V1ApiProps) {
	r.Route("/api/v1", v1Routes(p))
}
