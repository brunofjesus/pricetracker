package product

import (
	"strconv"
)

func filterEANs(storeProduct MqStoreProduct) []int64 {
	var validEans []int64
	for _, ean := range storeProduct.EAN {
		if eanInt, err := strconv.Atoi(ean); err == nil {
			validEans = append(validEans, int64(eanInt))
		}
	}

	return validEans
}
