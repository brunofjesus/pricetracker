package mq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brunofjesus/pricetracker/catalog/internal/app"
	"log/slog"

	"github.com/brunofjesus/pricetracker/catalog/pkg/product"
	"github.com/brunofjesus/pricetracker/catalog/pkg/store"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	productHandler           *product.Handler
	storeHandler             *store.Handler
	applicationConfiguration *app.ApplicationConfiguration
}

// Listen implements ProductConsumer.
func (c *Consumer) Listen(ctx context.Context) error {
	logger := ctx.Value("logger").(*slog.Logger)

	logger.Debug("connecting to MQ")
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

	logger.Debug("setup channel")
	err = channelSetup(ctx, ch)
	if err != nil {
		return err
	}

	manualAck := c.applicationConfiguration.MessageQueue.ManualAck

	logger.Debug("start consuming mq", slog.String("queue", "catalog"), slog.Bool("manualAck", manualAck))
	msgs, err := ch.Consume(
		"catalog",  // queue
		"",         // Consumer
		!manualAck, // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	if err != nil {
		return fmt.Errorf("error registering RabbitMQ Consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			var err error = nil
			payloadType := d.RoutingKey

			switch payloadType {
			case "store":
				var store store.MqStore
				err = json.Unmarshal(d.Body, &store)
				if err != nil {
					logger.Error("cannot unmarshall store", slog.Any("error", err))
					continue
				}
				err = c.storeHandler.Handle(ctx, store)
			case "product":
				var storeProduct product.MqStoreProduct
				err = json.Unmarshal(d.Body, &storeProduct)
				if err != nil {
					logger.Error("cannot unmarshall store product", slog.Any("error", err))
					continue
				}
				err = c.productHandler.Handle(ctx, storeProduct)
			}

			if manualAck {
				if err == nil {
					err = d.Acknowledger.Ack(d.DeliveryTag, false)
				} else {
					err = d.Acknowledger.Nack(d.DeliveryTag, false, false)
				}
				if err != nil {
					logger.Error("cannot send ack/nack to mq", slog.Any("error", err))
				}
			}
		}
	}()

	logger.Info("connected to mq")

	<-conn.NotifyClose(make(chan *amqp.Error))

	return errors.New("MQ Connection closed")
}

func (c *Consumer) connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial(c.applicationConfiguration.MessageQueue.URL)

	if err != nil {
		return nil, fmt.Errorf("error connecting to RabbitMQ: %w", err)
	}

	return conn, nil
}

func channelSetup(ctx context.Context, ch *amqp.Channel) error {
	logger := ctx.Value("logger").(*slog.Logger)

	// Prefetch 10
	err := ch.Qos(10, 0, false)
	if err != nil {
		return fmt.Errorf("error setting prefetch: %w", err)
	}

	exchangeName := "catalog_ex"

	logger.Debug(
		"mq channel exchange declare",
		slog.Group(
			"exchange",
			slog.String("name", exchangeName),
		),
	)
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

	queueName := "catalog"

	logger.Debug(
		"mq channel queue declare",
		slog.Group(
			"queue",
			slog.String("name", exchangeName),
		),
	)
	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return fmt.Errorf("error declaring RabbitMQ queue: %w", err)
	}

	productRoutingKey := "product"

	logger.Debug(
		"mq channel queue bind",
		slog.Group(
			"bind",
			slog.String("queue", queueName),
			slog.String("exchange", exchangeName),
			slog.String("routing_key", productRoutingKey),
		),
	)
	err = ch.QueueBind(
		queueName,         // queue name
		productRoutingKey, // routing key
		exchangeName,      // exchange name,
		false,             // no-wait
		nil,               // arguments
	)

	if err != nil {
		return fmt.Errorf(
			"error binding queue '%s' to exchange '%s' on routing key '%s': %v",
			"catalog", "catalog_ex", "product", err)
	}

	storeRoutingKey := "store"

	logger.Debug(
		"mq channel queue bind",
		slog.Group(
			"bind",
			slog.String("queue", queueName),
			slog.String("exchange", exchangeName),
			slog.String("routing_key", productRoutingKey),
		),
	)
	err = ch.QueueBind(
		queueName,       // queue name
		storeRoutingKey, // routing key
		exchangeName,    // exchange name,
		false,           // no-wait
		nil,             // arguments
	)

	if err != nil {
		return fmt.Errorf(
			"error binding queue '%s' to exchange '%s' on routing key '%s': %v",
			"catalog", "catalog_ex", "store", err)
	}

	return nil
}
