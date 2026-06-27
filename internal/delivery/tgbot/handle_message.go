package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleMessage(msg *tgbotapi.Message) {
	// Если это команда (начинается со слэша /)
	if msg.IsCommand() {
		b.handleCommand(msg)
		return
	}

	// Если это обычный текст (создание задачи)
	b.handleTaskCreation(msg)
}
