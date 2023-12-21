package frontend

import (
	"embed"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/frontend/pages"
	"github.com/brunofjesus/pricetracker/catalog/pkg/http/rest/utils"
	"github.com/go-chi/chi/v5"
	"io/fs"
	"log"
	"net/http"
)

type V1FrontendProps struct {
}

//go:embed static
var static embed.FS

func AddRoutes(r chi.Router, p V1FrontendProps) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		component := pages.Hello("bruno")
		err := component.Render(r.Context(), w)
		if err != nil {
			utils.ErrorJSON(w, err, 500)
		}
	})

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
