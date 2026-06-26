package tgbot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"proj/internal/entity"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleTaskCreation(msg *tgbotapi.Message) {
	// Создаем контекст с таймаутом на 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Обязательно освобождаем ресурсы!

	taskTitle := strings.TrimSpace(msg.Text)
	// Передаем этот контекст дальше по всем слоям до самой БД!
	user := &entity.UserTg{
		ID:       msg.From.ID,
		Username: msg.From.UserName,
	}
	err := b.userUC.GetUserByTgID(ctx, user)
	if err != nil {
		b.userUC.RegisterUserTg(ctx, user)
		b.send(msg.Chat.ID, "Привет! Теперь ты в системе. Отправь мне текст, чтобы создать задачу.")
		return
	}
	task, err := b.taskUC.CreateTask(ctx, taskTitle, user.UserID)
	if err != nil {
		// Если контекст завершился по таймауту
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			b.logger.Error("database timeout during task creating")
			b.send(msg.Chat.ID, "⚠️ База данных не ответила вовремя. Попробуйте еще раз.")
			return
		}
		b.send(msg.Chat.ID, fmt.Sprintf("⚠️ Ошибка: %s", err.Error()))
		b.logger.Error("failed to create task", slog.String("error", err.Error()))
		return
	}

	b.send(msg.Chat.ID, fmt.Sprintf("✅ Задача №%d создана!", task.ID))
}
