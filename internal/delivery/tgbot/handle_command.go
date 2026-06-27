package tgbot

import (
	"context"

	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleCommand(msg *tgbotapi.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch msg.Command() {
	case "start":
		b.handleStartCommand(ctx, msg)
	case "help":
		b.send(msg.Chat.ID, "Доступные команды:\n/start - Регистрация\n/list - Мои задачи")
	default:
		b.send(msg.Chat.ID, "Неизвестная команда 🤔")
	}
}
