package pingodoce

import (
	"fmt"

	"github.com/brunofjesus/pricetracker/store"
	"github.com/shopspring/decimal"
)

func mapPingoDoceProductToStoreProduct(in PingoDoceProduct) store.StoreProduct {
	return store.StoreProduct{
		StoreId:   "pingodoce",
		StoreName: "Pingo Doce",
		Name:      in.ShortDescription,
		EAN:       in.Eans,
		SKU:       []string{in.Sku},
		Brand:     in.Brand.Name,
		Price:     decimal.NewFromFloat(in.BuyingPrice),
		Available: in.OnlineStatus == "AVAILABLE",
		ImageLink: fmt.Sprintf("https://res.cloudinary.com/fonte-online/image/upload/c_fill,h_600,q_auto,w_600/v1/PDO_PROD/%s_1", in.Sku),
		Link:      fmt.Sprintf("https://mercadao.pt/store/pingo-doce/product/%s", in.Slug),
	}
}
