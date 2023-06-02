package worker

import (
	"fmt"
	"log"
	"strings"
	"time"

	"worker/internal/app/config"
	"worker/internal/app/rabbitmq"
	"worker/internal/app/telegram"
)

type Controller struct {
	r *rabbitmq.RabbitMq
	b *telegram.Bot
}

func StartWorker() error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}

	rabbit, err := rabbitmq.StartRabbitMq(conf.RabbitMQ)
	if err != nil {
		return err
	}

	defer func() {
		_ = rabbit.Conn.Close()
	}()

	bot, err := telegram.NewBot(conf)
	if err != nil {
		return err
	}

	c := &Controller{
		r: rabbit,
		b: bot,
	}

	var ch chan error

	c.worker(ch)

	for err = range ch {
		log.Print(err)
	}

	log.Print("service stopped")

	return nil
}

func (c *Controller) worker(ch chan<- error) {
	go func() {
		log.Print("starting goroutine")

		defer func() {
			c.worker(ch)
		}()

		for {
			get, err := c.r.GetMessage("maya")
			if err != nil {
				go func() {
					ch <- err
				}()
				continue
			}

			if get.Type != "" {
				fmt.Printf("[%s] Уведомление для %s от %s: %s\n",
					strings.ToUpper(get.Type),
					get.ReplyTo,
					get.Timestamp.Format(time.RFC822),
					get.Body)
			}

			switch strings.ToLower(get.Type) {
			case "telegram":
				if err := c.b.SendMessage(get); err != nil {
					if err := c.r.ReturnMessage(get); err != nil {
						go func() {
							ch <- err
						}()
						continue
					}
				}
			}
		}
	}()
}
