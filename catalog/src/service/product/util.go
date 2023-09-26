package product

import (
	"strconv"

	"github.com/brunofjesus/pricetracker/catalog/src/integration"
)

func filterEANs(storeProduct integration.StoreProduct) []int64 {
	var validEans []int64
	for _, ean := range storeProduct.EAN {
		if eanInt, err := strconv.Atoi(ean); err == nil {
			validEans = append(validEans, int64(eanInt))
		}
	}

	return validEans
}
