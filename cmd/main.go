package main

import (
	"fmt"

	"github.com/brunofjesus/pricetracker/store"
	"github.com/brunofjesus/pricetracker/store/pingodoce"
)

func main() {
	fmt.Println("Hello world")

	productChannel := make(chan store.StoreProduct, 10)

	var store = pingodoce.Instance()
	go store.Crawl(productChannel)

	for storeProduct := range productChannel {
		fmt.Printf("%v\n", storeProduct)
	}
}
