package mq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/brunofjesus/pricetracker/stores/pingodoce/definition"
	"github.com/rabbitmq/amqp091-go"
)

func PublishProduct(ch *amqp091.Channel, product definition.StoreProduct) error {
	messageBytes, err := json.Marshal(product)

	if err != nil {
		return err
	}

	return publish(ch, "product", messageBytes)
}

func PublishStore(ch *amqp091.Channel, store definition.Store) error {
	messageBytes, err := json.Marshal(store)

	if err != nil {
		return err
	}

	return publish(ch, "store", messageBytes)
}

func publish(ch *amqp091.Channel, routingKey string, messageBytes []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(ctx,
		"catalog_ex", // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		})
}
