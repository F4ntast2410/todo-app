package tgbot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	message := tgbotapi.NewMessage(msg.Chat.ID, "")
	switch msg.Command() {
	case "start":
		b.handleStartCommand(ctx, msg)
	case "help":
		message.Text = "Доступные команды:\n/start - Регистрация\n/list - Мои задачи"
		b.Send(message)
	case "list":
		b.sendTaskList(ctx, msg)
	default:
		message.Text = "Неизвестная команда 🤔"
		b.Send(message)
	}
}
