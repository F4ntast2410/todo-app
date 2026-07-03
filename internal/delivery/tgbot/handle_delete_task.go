package tgbot

import (
	"context"
	"errors"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handlerDeleteTask(ctx context.Context, query *tgbotapi.CallbackQuery, taskID int) {
	err := b.taskUC.DeleteTask(ctx, taskID)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			b.logger.Error("database timeout during task creating")
			message := tgbotapi.NewMessage(query.From.ID, "⚠️ База данных не ответила вовремя. Попробуйте еще раз.")
			b.Send(message)
			return
		}
		b.logger.Error("error deleting task", slog.String("error", err.Error()))
		return
	}
	b.editTaskViewCallback(ctx, query, taskID)
}

func (b *BotServer) handlerRecoverTask(ctx context.Context, query *tgbotapi.CallbackQuery, taskID int) {
	err := b.taskUC.RecoverTask(ctx, taskID)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			b.logger.Error("database timeout during task creating")
			message := tgbotapi.NewMessage(query.From.ID, "⚠️ База данных не ответила вовремя. Попробуйте еще раз.")
			b.Send(message)
			return
		}
		b.logger.Error("error recovering task", slog.String("error", err.Error()))
		return
	}
	b.editTaskViewCallback(ctx, query, taskID)
}
