package tgbot

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleMessage(msg *tgbotapi.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	message := tgbotapi.NewMessage(msg.Chat.ID, "")
	if msg.IsCommand() {
		b.handleCommand(ctx, msg)
		return
	}

	// Если это обычный текст (создание задачи)
	b.sessionCache.Mu.RLock()
	state := b.sessionCache.Cache[msg.From.ID].State
	b.sessionCache.Mu.RUnlock()

	switch state {
	case StateWaitingTaskTitle:
		b.handleGetTitleMessage(msg)
	case StateWaitingTaskDescription:
		b.handleTaskCreation(ctx, msg)
	case StateWaitingNewTaskDescription:
		b.handlerUpdateDecription(ctx, msg)
	default:
		message.Text = "Неизвестная команда 🤔"
		b.Send(message)
	}
}
