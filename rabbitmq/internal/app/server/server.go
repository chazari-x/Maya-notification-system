package server

import (
	"net/http"

	"Maya-notification-system/rabbitmq/internal/app/config"
	"Maya-notification-system/rabbitmq/internal/app/handler"
	"Maya-notification-system/rabbitmq/internal/app/rabbitmq"
	"github.com/go-chi/chi/v5"
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

	return http.ListenAndServe(conf.RunAddress, r)
}
