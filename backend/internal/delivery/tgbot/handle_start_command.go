package tgbot

import (
	"context"
	"errors"
	"log/slog"
	customErrors "proj/internal/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleStartCommand(ctx context.Context, msg *tgbotapi.Message) {
	message := tgbotapi.NewMessage(msg.Chat.ID, "Привет! Теперь ты в системе. Отправь мне текст, чтобы создать задачу.")
	err := b.userUC.RegisterUserTg(ctx, msg.From.ID, msg.From.UserName)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			b.logger.Error("database timeout during user registration")
			message.Text = "⏱️ Превышено время ожидания базы данных. Попробуйте позже."
			b.Send(message)
			return
		}
		if errors.Is(err, customErrors.ErrUserAlreadyExists) {
			message.Text = "Рад видеть тебя снова! Твои задачи на месте."
			b.Send(message)
			return
		}
		b.logger.Error("failed to register user", slog.String("error", err.Error()))
		message.Text = "⚠️ Не удалось завершить регистрацию."
		b.Send(message)
		return
	}

	b.Send(message)
}
