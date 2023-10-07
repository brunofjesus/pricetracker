package mq

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/config"
	"github.com/brunofjesus/pricetracker/catalog/src/integration"
	"github.com/brunofjesus/pricetracker/catalog/src/service/product"
	"github.com/brunofjesus/pricetracker/catalog/src/service/store"
	amqp "github.com/rabbitmq/amqp091-go"
)

var once sync.Once
var instance Consumer

type Consumer interface {
	Listen() error
}

type consumer struct {
	productHandler           product.ProductHandler
	storeHandler             store.StoreHandler
	applicationConfiguration *config.ApplicationConfiguration
}

func GetListener() Consumer {
	once.Do(func() {
		instance = &consumer{
			productHandler:           product.GetProductHandler(),
			storeHandler:             store.GetStoreHandler(),
			applicationConfiguration: config.GetApplicationConfiguration(),
		}
	})

	return instance
}

// Listen implements ProductConsumer.
func (c *consumer) Listen() error {
	conn, err := amqp.Dial(c.applicationConfiguration.MessageQueue.URL)

	if err != nil {
		return fmt.Errorf("error connecting to RabbitMQ: %w", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		return fmt.Errorf("error opening RabbitMQ channel: %w", err)
	}
	defer ch.Close()

	ch.Qos(10, 0, false)

	err = ch.ExchangeDeclare(
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

	q, err := ch.QueueDeclare(
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

	manualAck := c.applicationConfiguration.MessageQueue.ManualAck

	msgs, err := ch.Consume(
		q.Name,     // queue
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

	var forever chan struct{}

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
	<-forever

	return nil
}
