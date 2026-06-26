package tgbot

import (
	"context"
	"errors"
	"log/slog"
	"proj/internal/entity"
	"proj/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleStartCommand(ctx context.Context, msg *tgbotapi.Message) {

	user := entity.UserTg{
		ID:       msg.From.ID,
		Username: msg.From.UserName,
	}

	err := b.userUC.RegisterUserTg(ctx, &user)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			b.logger.Error("database timeout during user registration")
			b.send(msg.Chat.ID, "⏱️ Превышено время ожидания базы данных. Попробуйте позже.")
			return
		}
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			b.send(msg.Chat.ID, "Рад видеть тебя снова! Твои задачи на месте.")
			return
		}
		b.logger.Error("failed to register user", slog.String("error", err.Error()))
		b.send(msg.Chat.ID, "⚠️ Не удалось завершить регистрацию.")
		return
	}

	b.send(msg.Chat.ID, "Привет! Теперь ты в системе. Отправь мне текст, чтобы создать задачу.")
}
