package tgbot

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handlerDeleteTask(ctx context.Context, query *tgbotapi.CallbackQuery, taskID int) {
	err := b.taskUC.DeleteTask(ctx, taskID)
	if err != nil {
		b.logger.Error("error deleting task", slog.String("error", err.Error()))
		return
	}
	b.editTaskViewCallback(ctx, query, taskID)
}

func (b *BotServer) handlerRecoverTask(ctx context.Context, query *tgbotapi.CallbackQuery, taskID int) {
	err := b.taskUC.RecoverTask(ctx, taskID)
	if err != nil {
		b.logger.Error("error recovering task", slog.String("error", err.Error()))
		return
	}
	b.editTaskViewCallback(ctx, query, taskID)
}
