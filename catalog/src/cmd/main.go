package main

import (
	"fmt"
	"log"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"
	"github.com/brunofjesus/pricetracker/catalog/src/datasource/pingodoce"
	store_repository "github.com/brunofjesus/pricetracker/catalog/src/repository/store"
	"github.com/brunofjesus/pricetracker/catalog/src/service/product"
	store_service "github.com/brunofjesus/pricetracker/catalog/src/service/store"
)

func main() {
	fmt.Println("Hello world")

	storeRepository := store_repository.NewStoreRepository()
	storeEnroller := store_service.NewStoreEnroller(storeRepository)

	listener := product.GetListener()

	stores := listStores()
	for _, store := range stores {
		if err := storeEnroller.Enroll(store); err != nil {
			log.Printf("Cannot enroll %s: %v\n", store.Name(), err)
		} else {
			go store.Crawl(listener.ProductChannel)
		}
	}

	listener.Start()

	for {
		time.Sleep(1 * time.Second)
	}
}

func listStores() []datasource.Store {
	return []datasource.Store{
		pingodoce.Instance(),
	}
}
