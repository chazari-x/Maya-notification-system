package rabbitmq

import (
	"context"
	"log"
	"time"

	"Maya-notification-system/rabbitmq/internal/app/model"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	Conn *amqp091.Connection
}

func StartRabbitMq(rabbitmq string) (RabbitMq, error) {
	conn, err := amqp091.Dial(rabbitmq)
	if err != nil {
		return RabbitMq{}, err
	}

	return RabbitMq{Conn: conn}, nil
}

func (r *RabbitMq) SendMessage(msg model.MsgStruct) error {
	ch, err := r.Conn.Channel()
	if err != nil {
		return err
	}

	defer func() {
		_ = ch.Close()
	}()

	msgTime := time.Now().Format(time.RFC3339)

	q, err := ch.QueueDeclare(
		msg.MsgType,
		false,
		false,
		false,
		false,
		amqp091.Table{
			amqp091.ExchangeHeaders: msgTime,
			amqp091.ExchangeDirect:  msg.MsgTo,
		},
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg.Msg),
		}); err != nil {
		return err
	}

	log.Printf(" [x] Type: %s Sent %s\n", msg.MsgType, msg.Msg)
	return nil
}
