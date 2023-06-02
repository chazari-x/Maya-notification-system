package telegram

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rabbitmq/amqp091-go"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"worker/internal/app/config"
)

type Bot struct {
	bot  *tgbotapi.BotAPI
	conf *config.Config
}

func NewBot(config *config.Config) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("new bot api (%s) err: %s", config.Bot.Token, err.Error())
	}

	return &Bot{bot: bot, conf: config}, nil
}

func (b *Bot) SendMessage(msg amqp091.Delivery) error {
	id, err := strconv.Atoi(msg.ReplyTo)
	if err != nil {
		return err
	}

	text := fmt.Sprintf("Уведомление от %s:\n%s", msg.Timestamp.Format(time.RFC822), msg.Body)
	message := tgbotapi.NewMessage(int64(id), text)
	if _, err = b.bot.Send(message); err != nil {
		return err
	}

	return nil
}
