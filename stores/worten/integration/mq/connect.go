package mq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func Connect(url string) (*amqp091.Connection, *amqp091.Channel, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		defer conn.Close() // ups, close the channel
		return nil, nil, err
	}

	err = setup(ch)
	if err != nil {
		defer conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}

func setup(ch *amqp091.Channel) error {
	_, err := ch.QueueDeclare(
		"catalog", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return fmt.Errorf("error declaring RabbitMQ queue: %w", err)
	}

	err = ch.QueueBind(
		"catalog",    // queue name
		"product",    // routing key
		"catalog_ex", // exchange name,
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return fmt.Errorf(
			"error binding queue '%s' to exchange '%s' on routing key '%s': %v",
			"catalog", "catalog_ex", "product", err)
	}

	err = ch.QueueBind(
		"catalog",    // queue name
		"store",      // routing key
		"catalog_ex", // exchange name,
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return fmt.Errorf(
			"error binding queue '%s' to exchange '%s' on routing key '%s': %v",
			"catalog", "catalog_ex", "store", err)
	}

	return nil
}
