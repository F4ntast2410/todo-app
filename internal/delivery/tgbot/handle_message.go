package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleMessage(msg *tgbotapi.Message) {
	// Если это команда (начинается со слэша /)
	message := tgbotapi.NewMessage(msg.Chat.ID, "")
	if msg.IsCommand() {
		b.handleCommand(msg)
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
		b.handleTaskCreation(msg)
	default:
		message.Text = "Неизвестная команда 🤔"
		b.Send(message)
	}
}
