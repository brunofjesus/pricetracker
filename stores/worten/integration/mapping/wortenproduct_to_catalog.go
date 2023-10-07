package mapping

import (
	"fmt"

	"github.com/brunofjesus/pricetracker/stores/worten/definition/catalog"
	"github.com/brunofjesus/pricetracker/stores/worten/definition/store"
	"github.com/shopspring/decimal"
)

func MapWortenProductToCatalogProduct(source store.WortenProductHit, destination *catalog.StoreProduct) error {
	winningOffer, err := decimal.NewFromString(source.WinningOffer.Price.Value)
	if err != nil {
		return fmt.Errorf("cannot handle winning offer (%v) %v", source.WinningOffer.Price, err)
	}

	var price = winningOffer.Div(decimal.NewFromInt(100))

	if source.SecondOfferPrice != nil && source.SecondOfferPrice.Value != "0" {
		secondOffer, err := decimal.NewFromString(source.SecondOfferPrice.Value)
		if err != nil {
			return fmt.Errorf("cannot handle second offer (%v) %v", source.WinningOffer.Price, err)
		}

		price = decimal.Min(winningOffer, secondOffer).Div(decimal.NewFromInt(100))
	}

	destination.StoreSlug = "worten"
	destination.EAN = source.Product.Ean
	destination.SKU = []string{source.Product.Sku}
	destination.Name = source.Product.Name
	destination.Brand = source.Product.BrandName
	destination.Price = price
	destination.Available = source.WinningOffer.IsInStock
	destination.ImageLink = source.Product.Image.URL
	destination.Link = fmt.Sprintf("https://www.worten.pt/%s", source.Product.URL)

	return nil
}
