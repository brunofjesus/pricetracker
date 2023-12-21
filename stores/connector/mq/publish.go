package mq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brunofjesus/pricetracker/stores/connector/dto"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"
)

type Publisher struct {
	logger          *slog.Logger
	messageQueueUrl string
	conn            *amqp091.Connection
	ch              *amqp091.Channel
	closeChan       chan struct{}
}

func NewPublisher(logger *slog.Logger, messageQueueUrl string) (*Publisher, error) {
	conn, ch, err := Connect(messageQueueUrl)
	if err != nil {
		return nil, err
	}

	instance := &Publisher{
		logger:          logger,
		messageQueueUrl: messageQueueUrl,
		conn:            conn,
		ch:              ch,
		closeChan:       make(chan struct{}),
	}

	go instance.backgroundLoop()

	return instance, nil
}

func (p *Publisher) PublishProduct(product dto.StoreProduct) error {
	messageBytes, err := json.Marshal(product)

	if err != nil {
		return err
	}

	return p.publish("product", messageBytes)
}

func (p *Publisher) PublishStore(store dto.Store) error {
	messageBytes, err := json.Marshal(store)

	if err != nil {
		return err
	}

	return p.publish("store", messageBytes)
}

func (p *Publisher) Close() error {
	err := errors.Join(
		p.conn.Close(),
		p.ch.Close(),
	)

	p.closeChan <- struct{}{}

	return err
}

func (p *Publisher) publish(routingKey string, messageBytes []byte) error {

	for p.conn == nil || p.ch == nil || p.conn.IsClosed() || p.ch.IsClosed() {
		p.logger.Info("Connection or channel are not set, waiting...")
		time.Sleep(1 * time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.ch.PublishWithContext(ctx,
		"catalog_ex", // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		})
}

func (p *Publisher) backgroundLoop() {
	for {
		select {
		case amqpErr := <-p.conn.NotifyClose(make(chan *amqp091.Error)):
			fmt.Print(amqpErr)
			p.logger.Info("MQ Notify close event")
			p.handleMqClose(amqpErr)
			p.logger.Debug("MQ Notify close channel terminated")
		case <-p.closeChan:
			p.logger.Info("Stopping publisher background loop!")
			return
		}
	}
}

func (p *Publisher) handleMqClose(amqpErr *amqp091.Error) {
	if !p.conn.IsClosed() {
		_ = p.conn.Close()
	}
	if !p.ch.IsClosed() {
		_ = p.ch.Close()
	}

	//it was an error
	if amqpErr != nil {
		for {
			conn, ch, err := Connect(p.messageQueueUrl)
			if err != nil {
				p.logger.Error("MQ Reconnect error", slog.Any("error", err))
				time.Sleep(1 * time.Second)
			} else {
				p.logger.Debug("Re connected")
				p.conn = conn
				p.ch = ch
				break
			}
		}
	}
}
