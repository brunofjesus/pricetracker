package mapping

import (
	"fmt"
	"github.com/brunofjesus/pricetracker/stores/connector/dto"
	"github.com/brunofjesus/pricetracker/stores/worten/pkg/definition"
	"strconv"
)

func MapWortenProductToCatalogProduct(source definition.WortenProductHit, destination *dto.StoreProduct) error {
	winningOffer, err := strconv.ParseInt(source.WinningOffer.Price.Value, 10, 64)
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
	destination.Currency = "EUR"

	return nil
}
