package main

import (
	"fmt"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/management"
	"github.com/brunofjesus/pricetracker/catalog/src/store/pingodoce"
)

func main() {
	//TODO
	// https://github.com/glebarez/go-sqlite
	// https://github.com/golang-migrate/migrate/tree/master/database/sqlite

	fmt.Println("Hello world")

	listener := management.GetListener()

	var store = pingodoce.Instance()
	go store.Crawl(listener.ProductChannel)

	listener.Start()

	for {
		time.Sleep(1 * time.Second)
	}
}
