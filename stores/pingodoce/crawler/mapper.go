package crawler

import (
	"fmt"

	"github.com/brunofjesus/pricetracker/stores/pingodoce/definition"
)

func mapPingoDoceProductToStoreProduct(store definition.Store, in definition.PingoDoceProduct) definition.StoreProduct {
	return definition.StoreProduct{
		StoreSlug: store.Slug,
		Name:      in.ShortDescription,
		EAN:       in.Eans,
		SKU:       []string{in.Sku},
		Brand:     in.Brand.Name,
		Price:     int64(in.BuyingPrice * 100),
		Available: in.OnlineStatus == "AVAILABLE",
		ImageLink: fmt.Sprintf("https://res.cloudinary.com/fonte-online/image/upload/c_fill,h_600,q_auto,w_600/v1/PDO_PROD/%s_1", in.Sku),
		Link:      fmt.Sprintf("https://mercadao.pt/store/pingo-doce/product/%s", in.Slug),
	}
}
