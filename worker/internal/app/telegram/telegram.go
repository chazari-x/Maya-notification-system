package telegram

import (
	"fmt"
	"strconv"

	"Maya-notification-system/worker/internal/app/config"
	"github.com/rabbitmq/amqp091-go"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Bot struct {
	bot  *tgbotapi.BotAPI
	conf *config.Config
}

func NewBot(config *config.Config) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("new bot api err: %s", err.Error())
	}

	return &Bot{bot: bot, conf: config}, nil
}

func (b *Bot) SendMessage(msg amqp091.Delivery) error {
	id, err := strconv.Atoi(msg.ReplyTo)
	if err != nil {
		return err
	}

	text := fmt.Sprintf("Уведомление от %s:\n%s", msg.Timestamp, msg.Body)
	message := tgbotapi.NewMessage(int64(id), text)
	if _, err = b.bot.Send(message); err != nil {
		return err
	}

	return nil
}
