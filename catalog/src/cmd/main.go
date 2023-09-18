package main

import (
	"fmt"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/management"
	"github.com/brunofjesus/pricetracker/catalog/src/store/pingodoce"
)

func main() {
	fmt.Println("Hello world")

	dsn := "host=localhost port=5432 user=postgres password=price dbname=postgres sslmode=disable timezone=UTC connect_timeout=5"

	db := databaseConnect(dsn, 5)

	fmt.Println(db.Stats().OpenConnections)

	listener := management.GetListener()

	var store = pingodoce.Instance()
	go store.Crawl(listener.ProductChannel)

	listener.Start()

	for {
		time.Sleep(1 * time.Second)
	}
}
