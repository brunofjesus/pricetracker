package mq

import "github.com/rabbitmq/amqp091-go"

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

	return conn, ch, nil
}
