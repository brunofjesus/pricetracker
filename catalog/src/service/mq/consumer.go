package mq

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/brunofjesus/pricetracker/catalog/src/config"
	"github.com/brunofjesus/pricetracker/catalog/src/integration"
	"github.com/brunofjesus/pricetracker/catalog/src/service/product"
	"github.com/brunofjesus/pricetracker/catalog/src/service/store"
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Listen() error
}

type consumer struct {
	productHandler           product.ProductHandler
	storeHandler             store.StoreHandler
	applicationConfiguration *config.ApplicationConfiguration
}

func NewConsumer() Consumer {
	return &consumer{
		productHandler:           product.GetProductHandler(),
		storeHandler:             store.GetStoreHandler(),
		applicationConfiguration: config.GetApplicationConfiguration(),
	}
}

// Listen implements ProductConsumer.
func (c *consumer) Listen() error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		return fmt.Errorf("error opening RabbitMQ channel: %w", err)
	}
	defer ch.Close()

	err = channelSetup(ch)
	if err != nil {
		return err
	}

	manualAck := c.applicationConfiguration.MessageQueue.ManualAck

	msgs, err := ch.Consume(
		"catalog",  // queue
		"",         // consumer
		!manualAck, // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	if err != nil {
		return fmt.Errorf("error registering RabbitMQ consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			var err error = nil
			payloadType := d.RoutingKey

			switch payloadType {
			case "store":
				var store integration.Store
				err = json.Unmarshal(d.Body, &store)
				if err != nil {
					log.Printf("error unmarshalling store: %v", err)
					continue
				}
				err = c.storeHandler.Handle(store)
			case "product":
				var storeProduct integration.StoreProduct
				err = json.Unmarshal(d.Body, &storeProduct)
				if err != nil {
					log.Printf("error unmarshalling store product: %v", err)
					continue
				}
				err = c.productHandler.Handle(storeProduct)
			}

			if manualAck {
				if err == nil {
					err = d.Acknowledger.Ack(d.DeliveryTag, false)
				} else {
					err = d.Acknowledger.Nack(d.DeliveryTag, false, false)
				}
				if err != nil {
					log.Printf("cannot send ack/nack to mq: %v", err)
				}
			}
		}
	}()

	log.Printf("Connected to Message Queue.")

	<-conn.NotifyClose(make(chan *amqp.Error))

	return errors.New("MQ Connection closed")
}

func (c *consumer) connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial(c.applicationConfiguration.MessageQueue.URL)

	if err != nil {
		return nil, fmt.Errorf("error connecting to RabbitMQ: %w", err)
	}

	return conn, nil
}

func channelSetup(ch *amqp091.Channel) error {
	// Prefetch 10
	ch.Qos(10, 0, false)

	err := ch.ExchangeDeclare(
		"catalog_ex", // name
		"direct",     //kind
		false,        // durable
		false,        // auto delete,
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return fmt.Errorf("error declaring exchange: %w", err)
	}

	_, err = ch.QueueDeclare(
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
