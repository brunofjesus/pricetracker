package rest

import (
	"fmt"
	"github.com/brunofjesus/pricetracker/catalog/pkg/price"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"github.com/brunofjesus/pricetracker/catalog/pkg/rest/handler"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type PropsV1 struct {
	ProductFinder *product.Finder
	PriceFinder   *price.Finder
}

type RouteFunc = func(r chi.Router)

func ListenAndServe(propsV1 PropsV1, port int) {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", v1Routes(propsV1))

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func v1Routes(p PropsV1) RouteFunc {
	return func(r chi.Router) {
		r.Get("/product/{productId}", handler.GetProduct(p.ProductFinder))
		r.Get("/product/search", handler.SearchProduct(p.ProductFinder))
		r.Get("/product/{productId}/history", handler.GetHistory(p.PriceFinder))
	}
}
