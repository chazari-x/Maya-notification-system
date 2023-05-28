package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"Maya-notification-system/rabbitmq/internal/app/config"
	"Maya-notification-system/rabbitmq/internal/app/model"
	"Maya-notification-system/rabbitmq/internal/app/rabbitmq"
)

type Controller struct {
	c      config.Config
	rabbit rabbitmq.RabbitMq
}

func NewController(c config.Config, rabbit rabbitmq.RabbitMq) *Controller {
	return &Controller{c: c, rabbit: rabbit}
}

func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print("Post: read all err: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if string(b) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := model.MsgStruct{}
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Print("Post: json unmarshal err: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.rabbit.SendMessage(msg)
	if err != nil {
		log.Printf("Post: %s, msgType: %s, msg: %s",
			err.Error(), msg.MsgType, msg.Msg)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Post: %d, msgType: %s, msg: %s",
		http.StatusOK, msg.MsgType, msg.Msg)
	w.WriteHeader(http.StatusOK)
}
