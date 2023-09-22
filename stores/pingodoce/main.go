package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Pingo Doce Store Crawler 1.0")

	conn, err := amqp.Dial("amqp://user:price@localhost:5672/")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"product", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		panic(err)
	}

	Crawl(func(storeProduct StoreProduct) {
		publish(ch, storeProduct)
	})
}

func publish(ch *amqp.Channel, product StoreProduct) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messageBytes, err := json.Marshal(product)

	if err != nil {
		return err
	}

	return ch.PublishWithContext(ctx,
		"",        // exchange
		"product", // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		})
}
