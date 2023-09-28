package main

import (
	"fmt"

	"github.com/brunofjesus/pricetracker/stores/worten/service"
)

func main() {
	categories, err := service.FindCategories()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", categories)
}
