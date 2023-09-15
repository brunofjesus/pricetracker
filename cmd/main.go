package main

import (
	"fmt"

	"github.com/brunofjesus/pricetracker/store/pingodoce"
)

func main() {
	fmt.Println("Hello world")

	var store = pingodoce.Instance()
	store.Crawl(nil)
}
