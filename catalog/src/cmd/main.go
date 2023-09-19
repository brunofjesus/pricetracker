package main

import (
	"log"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"
	"github.com/brunofjesus/pricetracker/catalog/src/datasource/pingodoce"
	product_listener "github.com/brunofjesus/pricetracker/catalog/src/service/product/listener"
	store_service "github.com/brunofjesus/pricetracker/catalog/src/service/store"
)

func main() {
	productListener := product_listener.GetListener()

	stores := listStores()
	storeEnroller := store_service.GetStoreEnroller()
	for _, store := range stores {
		if err := storeEnroller.Enroll(store); err != nil {
			log.Printf("Cannot enroll %s: %v\n", store.Name(), err)
		} else {
			go store.Crawl(productListener.Channel())
		}
	}

	productListener.Listen()

	for {
		time.Sleep(1 * time.Second)
	}
}

func listStores() []datasource.Store {
	return []datasource.Store{
		pingodoce.Instance(),
	}
}
