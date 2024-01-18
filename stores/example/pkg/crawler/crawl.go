package crawler

import (
	"github.com/brunofjesus/pricetracker/stores/connector/dto"
	"math/rand"
)

const StoreSlug = "example-store"

type CallbackFunction = func(product dto.StoreProduct)

// Crawl does not actually crawl for products, but in a real case scenario you might need to do that,
// instead this method gets a pre-defined list of products, applies random price variations and calls
// the provided callback function, which should publish the product to the catalog.
func Crawl(callback CallbackFunction) {
	for i, _ := range products {
		product := &products[i]
		randomizePrice(product)
		callback(*product)
	}

}

func randomizePrice(product *dto.StoreProduct) {
	// generate random number between -10000 and 10000
	minimum := -10000
	maximum := 10000
	priceChange := rand.Intn(maximum-minimum) + minimum

	product.Price += int64(priceChange)
}
