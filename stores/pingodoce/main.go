package main

import (
	"fmt"
	"log"

	"github.com/brunofjesus/pricetracker/stores/pingodoce/config"
	"github.com/brunofjesus/pricetracker/stores/pingodoce/crawler"
	"github.com/brunofjesus/pricetracker/stores/pingodoce/definition"
	"github.com/brunofjesus/pricetracker/stores/pingodoce/mq"
)

const (
	StoreSlug    = "pingo-doce"
	StoreName    = "Pingo Doce"
	StoreWebSite = "https://mercadao.pt/store/pingo-doce"
)

func main() {
	log.Printf("%s crawler\n", StoreName)

	appConfig := config.GetApplicationConfiguration()

	conn, ch, err := mq.Connect(appConfig.MessageQueue.URL)

	if err != nil {
		panic(err)
	}

	defer conn.Close()
	defer ch.Close()

	store := definition.Store{
		Slug:    StoreSlug,
		Name:    StoreName,
		Website: StoreWebSite,
	}

	err = mq.PublishStore(ch, store)

	if err != nil {
		log.Fatal(err)
	}

	crawler.Crawl(store, func(storeProduct definition.StoreProduct) {
		err := mq.PublishProduct(ch, storeProduct)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
	})

	log.Printf("Bye!")
}
