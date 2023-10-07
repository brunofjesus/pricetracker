package main

import (
	"log"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/config"
	"github.com/brunofjesus/pricetracker/catalog/src/service/mq"
)

func main() {
	appConfig := config.GetApplicationConfiguration()

	for i := 0; i < appConfig.MessageQueue.ThreadCount; i++ {
		go worker(i + 1)
	}

	select {}
}

func worker(id int) {
	for {
		log.Printf("(Re)Starting worker %d", id)
		err := mq.NewConsumer().Listen()
		if err != nil {
			log.Printf("Error on worker %d: %v", id, err)
		}
		time.Sleep(5000 * time.Millisecond)
	}
}
