package mapping

import (
	"fmt"
	"strconv"

	"github.com/brunofjesus/pricetracker/stores/worten/definition/catalog"
	"github.com/brunofjesus/pricetracker/stores/worten/definition/store"
)

func MapWortenProductToCatalogProduct(source store.WortenProductHit, destination *catalog.StoreProduct) error {
	winningOffer, err := strconv.Atoi(source.WinningOffer.Price.Value)
	if err != nil {
		return fmt.Errorf("cannot handle winning offer (%v) %v", source.WinningOffer.Price, err)
	}

	destination.StoreSlug = "worten"
	destination.EAN = source.Product.Ean
	destination.SKU = []string{source.Product.Sku}
	destination.Name = source.Product.Name
	destination.Brand = source.Product.BrandName
	destination.Price = winningOffer
	destination.Available = source.WinningOffer.IsInStock
	destination.ImageLink = source.Product.Image.URL
	destination.Link = fmt.Sprintf("https://www.worten.pt/%s", source.Product.URL)

	return nil
}
