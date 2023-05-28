package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"rabbitmq/internal/app/config"
	"rabbitmq/internal/app/handler"
	"rabbitmq/internal/app/rabbitmq"
)

func StartServer() error {
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

	c := handler.NewController(conf, rabbit)

	r := chi.NewRouter()
	r.Post("/api/send", c.Post)

	log.Print("starting server ", conf.RunAddress)
	return http.ListenAndServe(conf.RunAddress, r)
}
