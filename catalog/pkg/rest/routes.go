package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/pkg/rest/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ListenAndServe(port int) {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", v1Routes)

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func v1Routes(r chi.Router) {
	r.Get("/product/{productId}", handler.GetProduct)
	r.Get("/product/search", handler.SearchProduct)
	r.Get("/product/{productId}/history", handler.GetHistory)
}
