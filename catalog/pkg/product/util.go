package product

import (
	"strconv"
)

func filterNumbers(values []string) []int64 {
	var result []int64
	for _, value := range values {
		if number, err := strconv.Atoi(value); err == nil {
			result = append(result, int64(number))
		}
	}

	return result
}
