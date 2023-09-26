package main

import (
	"github.com/brunofjesus/pricetracker/catalog/src/service/mq"
)

func main() {

	listener := mq.GetListener()

	err := listener.Listen()

	if err != nil {
		panic(err)
	}
}
