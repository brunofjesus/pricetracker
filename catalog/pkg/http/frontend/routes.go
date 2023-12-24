package frontend

import (
	"embed"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/frontend/handler"
	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"github.com/brunofjesus/pricetracker/catalog/pkg/store"
	"github.com/go-chi/chi/v5"
	"io/fs"
	"log"
	"net/http"
)

type V1FrontendProps struct {
	ProductFinder *product.Finder
	StoreFinder   *store.Finder
}

//go:embed static
var static embed.FS

func AddRoutes(r chi.Router, p V1FrontendProps) {
	r.Get("/", handler.SearchProduct(p.ProductFinder, p.StoreFinder))
	r.Handle("/*", serveStatic())
}

func serveStatic() http.Handler {
	serverRoot, err := fs.Sub(static, "static")
	if err != nil {
		log.Print(err)
		return nil
	}
	var staticFS = http.FS(serverRoot)
	return http.FileServer(staticFS)
}
