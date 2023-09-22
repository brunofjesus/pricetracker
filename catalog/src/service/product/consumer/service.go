package consumer

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/integration"
	"github.com/brunofjesus/pricetracker/catalog/src/service/product/receiver"
	amqp "github.com/rabbitmq/amqp091-go"
)

var once sync.Once
var instance ProductConsumer

type ProductConsumer interface {
	Start()
}

type productConsumer struct {
	productReceiver receiver.ProductReceiver
}

func GetListener() ProductConsumer {
	once.Do(func() {
		instance = &productConsumer{
			productReceiver: receiver.GetProductReceiver(),
		}
	})

	return instance
}

// Start implements ProductConsumer.
func (c *productConsumer) Start() {
	conn, err := amqp.Dial("amqp://user:price@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"product", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			var storeProduct integration.StoreProduct
			err := json.Unmarshal(d.Body, &storeProduct)

			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			c.productReceiver.Receive(storeProduct)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
