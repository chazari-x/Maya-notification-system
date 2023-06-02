package worker

import (
	"fmt"
	"log"
	"strings"

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

			switch strings.ToLower(get.Type) {
			case "test":
				t := strings.Repeat("=", 55-len(get.Type)) +
					" [" + strings.ToUpper(get.Type) +
					"] " + strings.Repeat("=", 55-len(get.Type))

				fmt.Println(t)
				fmt.Printf("Сообщение для пользователя: %s\n", get.ReplyTo)
				fmt.Printf("От %s\n", get.Timestamp)
				fmt.Println(string(get.Body))
				fmt.Println(strings.Repeat("=", len(t)))
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
