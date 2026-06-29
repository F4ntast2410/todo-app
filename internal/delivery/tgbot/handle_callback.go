package tgbot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleCallback(query *tgbotapi.CallbackQuery) {
	data := strings.Split(query.Data, ":")

	switch data[0] {
	case "create_task_step":
		b.handleGetTitleCallback(query)
	}
}
