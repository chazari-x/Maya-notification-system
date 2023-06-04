package worker

import (
	"log"
	"strings"
	"sync"
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

	c.mainWorker()

	var wg sync.WaitGroup

	go func() {
		for true {
			if workerOne.work {
				wg.Add(1)
				c.worker(&wg, &workerOne)
			}

			wg.Wait()
		}
	}()

	for true {
		if workerTwo.work {
			wg.Add(1)
			c.worker(&wg, &workerTwo)
		}

		wg.Wait()
	}

	log.Print("service stopped")

	return nil
}

var (
	workerOne = workerStruct{work: false, nameNewWorker: "worker two"}
	workerTwo = workerStruct{work: false, nameNewWorker: "worker three"}
)

type workerStruct struct {
	work          bool
	nameNewWorker string
}

func (c *Controller) mainWorker() {
	go func() {
		log.Print("starting goroutine")

		defer func() {
			c.mainWorker()
		}()

		for {
			workerOne.work = false

			get, err := c.r.GetMessage("maya")
			if err != nil {
				log.Printf("worker get message err: %s", err)
				continue
			}

			if get.Body == nil {
				continue
			}

			workerOne.work = true

			if get.Type != "" {
				log.Printf("[%s] Уведомление для %s от %s: %s\n",
					strings.ToUpper(get.Type),
					get.ReplyTo,
					get.Timestamp.Format(time.RFC822),
					get.Body)
			}

			switch strings.ToLower(get.Type) {
			case "telegram":
				if err := c.b.SendMessage(get); err != nil {
					log.Printf("worker send message err: %s", err)
					if err := c.r.ReturnMessage(get); err != nil {
						log.Printf("worker return message err: %s", err)
						continue
					}
				}
			}
		}
	}()
}

func (c *Controller) worker(wg *sync.WaitGroup, w *workerStruct) {
	go func() {
		log.Printf("starting %s", w.nameNewWorker)

		for {
			if w.nameNewWorker == workerOne.nameNewWorker {
				w.work = false
			}

			if !w.work {
				log.Printf("%s is closed", w.nameNewWorker)
				wg.Done()
				break
			}

			get, err := c.r.GetMessage("maya")
			if err != nil {
				log.Printf("worker get message err: %s", err)
				continue
			}

			if get.Body == nil {
				continue
			}

			if w.nameNewWorker == workerOne.nameNewWorker {
				w.work = true
			}

			if get.Type != "" {
				log.Printf("[%s] Уведомление для %s от %s: %s\n",
					strings.ToUpper(get.Type),
					get.ReplyTo,
					get.Timestamp.Format(time.RFC822),
					get.Body)
			}

			switch strings.ToLower(get.Type) {
			case "telegram":
				if err := c.b.SendMessage(get); err != nil {
					log.Printf("worker send message err: %s", err)
					if err := c.r.ReturnMessage(get); err != nil {
						log.Printf("worker return message err: %s", err)
						continue
					}
				}
			}
		}
	}()
}
