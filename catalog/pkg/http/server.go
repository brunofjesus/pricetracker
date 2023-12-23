package http

import (
	"fmt"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/frontend"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/rest"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

type ServerProps struct {
	ApiProps      *rest.V1ApiProps
	FrontendProps *frontend.V1FrontendProps
	Host          string
	Port          int
}

func ListenAndServe(props ServerProps) error {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Use(middleware.Timeout(60 * time.Second))

	if props.ApiProps != nil {
		rest.AddRoutes(r, *props.ApiProps)
	}

	if props.FrontendProps != nil {
		frontend.AddRoutes(r, *props.FrontendProps)
	}

	return http.ListenAndServe(fmt.Sprintf("%s:%d", props.Host, props.Port), r)
}
