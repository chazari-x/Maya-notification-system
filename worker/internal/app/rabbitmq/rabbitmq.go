package rabbitmq

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
)

type RabbitMq struct {
	Conn *amqp091.Connection
}

func StartRabbitMq(rabbitmq string) (*RabbitMq, error) {
	conn, err := amqp091.Dial(rabbitmq)
	if err != nil {
		return nil, err
	}

	return &RabbitMq{Conn: conn}, nil
}

func (r *RabbitMq) GetMessage(msgType string) (amqp091.Delivery, error) {
	ch, err := r.Conn.Channel()
	if err != nil {
		return amqp091.Delivery{}, err
	}

	defer func() {
		_ = ch.Close()
	}()

	get, b, err := ch.Get(msgType, true)
	if err != nil {
		return amqp091.Delivery{}, err
	}

	if b {
		switch strings.ToLower(get.Type) {
		case "test", "telegram":
			return get, nil
		default:
			return amqp091.Delivery{}, ErrMethodNotAllowed
		}
	}

	return amqp091.Delivery{}, nil
}

func (r *RabbitMq) ReturnMessage(msg amqp091.Delivery) error {
	ch, err := r.Conn.Channel()
	if err != nil {
		return err
	}

	defer func() {
		_ = ch.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = ch.PublishWithContext(ctx,
		"",
		"maya",
		false,
		false,
		amqp091.Publishing{
			Body:        msg.Body,
			ReplyTo:     msg.ReplyTo,
			Type:        msg.Type,
			Timestamp:   msg.Timestamp,
			ContentType: "text/plain",
		}); err != nil {
		return err
	}

	log.Printf(" [x] msgType: %s, msg: %s, msgTo: %s\n", msg.Type, msg.Body, msg.ReplyTo)
	return nil
}
