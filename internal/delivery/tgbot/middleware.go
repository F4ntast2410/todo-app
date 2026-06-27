package tgbot

import (
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageHandler func(msg *tgbotapi.Message)

func WithLogging(next MessageHandler) MessageHandler {
	return MessageHandler(func(msg *tgbotapi.Message) {
		start := time.Now()

		next(msg)

		duration := time.Since(start)

		slog.Info("tg_bot message",
			slog.String("username", msg.From.UserName),
			slog.Int64("user_id", msg.Chat.ID),
			slog.String("text", msg.Text),
			slog.Duration("duration", duration),
		)
	})
}
