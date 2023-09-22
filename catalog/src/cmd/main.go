package main

import "github.com/brunofjesus/pricetracker/catalog/src/service/product/consumer"

func main() {

	productListener := consumer.GetListener()

	productListener.Start()

	// stores := listStores()
	// storeEnroller := store_service.GetStoreEnroller()
	// for _, store := range stores {
	// 	if err := storeEnroller.Enroll(store); err != nil {
	// 		log.Printf("Cannot enroll %s: %v\n", store.Name(), err)
	// 	} else {
	// 		go store.Crawl(productListener.Channel())
	// 	}
	// }

	// productListener.Listen()

	// for {
	// 	time.Sleep(1 * time.Second)
	// }
}
