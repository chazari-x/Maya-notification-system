package rabbitmq

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"Maya-notification-system/rabbitmq/internal/app/model"
	"github.com/rabbitmq/amqp091-go"
)

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
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
	switch strings.ToLower(msg.MsgType) {
	case "telegram", "test":
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
				Body:        []byte(msg.Msg),
				ReplyTo:     msg.MsgTo,
				Type:        msg.MsgType,
				Timestamp:   time.Now(),
				ContentType: "text/plain",
			}); err != nil {
			return err
		}

		log.Printf(" [x] msgType: %s, msg: %s, msgTo: %s\n", msg.MsgType, msg.Msg, msg.MsgTo)
	default:
		log.Printf(" [x] msgType: %s, msg: %s, msgTo: %s\n", msg.MsgType, msg.Msg, msg.MsgTo)
		return ErrMethodNotAllowed
	}

	return nil
}
