package main

import (
	"github.com/brunofjesus/pricetracker/catalog/src/config"
	"github.com/brunofjesus/pricetracker/catalog/src/service/mq"
)

func main() {
	appConfig := config.GetApplicationConfiguration()

	listener := mq.GetListener()

	for i := 0; i < appConfig.MessageQueue.ThreadCount; i++ {
		go listener.Listen()
	}

	select {}
}
