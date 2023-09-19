package main

import (
	"fmt"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/config"
	"github.com/brunofjesus/pricetracker/catalog/src/management"
	"github.com/brunofjesus/pricetracker/catalog/src/store/pingodoce"
)

func main() {
	fmt.Println("Hello world")

	appConfig := config.GetApplicationConfiguration()

	fmt.Printf("Connecting: %s\n", appConfig.Database.DSN)

	db := databaseConnect(appConfig.Database.DSN, 5)

	fmt.Println(db.Stats().OpenConnections)

	listener := management.GetListener()

	var store = pingodoce.Instance()
	go store.Crawl(listener.ProductChannel)

	listener.Start()

	for {
		time.Sleep(1 * time.Second)
	}
}
