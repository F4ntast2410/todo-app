package tgbot

import (
	"context"

	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleCommand(msg *tgbotapi.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	message := tgbotapi.NewMessage(msg.Chat.ID, "")
	switch msg.Command() {
	case "start":
		b.handleStartCommand(ctx, msg)
	case "help":
		message.Text = "Доступные команды:\n/start - Регистрация\n/list - Мои задачи"
		b.Send(message)
	case "list":
		b.handleTaskList(ctx, msg)
	default:
		message.Text = "Неизвестная команда 🤔"
		b.Send(message)
	}
}
