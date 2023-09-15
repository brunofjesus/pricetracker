package main

import (
	"fmt"
	"time"

	"github.com/brunofjesus/pricetracker/management"
	"github.com/brunofjesus/pricetracker/store/pingodoce"
)

func main() {
	fmt.Println("Hello world")

	listener := management.GetListener()

	var store = pingodoce.Instance()
	go store.Crawl(listener.ProductChannel)

	listener.Start()

	for {
		time.Sleep(1 * time.Second)
	}
}
