package app

import "github.com/brunofjesus/pricetracker/catalog/pkg/product"

type ApplicationEnvironment struct {
	Product struct {
		Handler product.ProductHandler
		Matcher product.ProductMatcher
	}
}
