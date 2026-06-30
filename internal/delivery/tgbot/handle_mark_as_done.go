package tgbot

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handlerMarkAsDone(ctx context.Context, query *tgbotapi.CallbackQuery, taskID int) {
	task, err := b.taskUC.GetTask(ctx, taskID)
	if err != nil {
		b.logger.Error("error getting task", slog.String("error: ", err.Error()))
		return
	}
	switch task.Done {
	case true:
		err = b.taskUC.MarkAsDone(ctx, taskID, false)
	case false:
		err = b.taskUC.MarkAsDone(ctx, taskID, true)
	}
	if err != nil {
		b.logger.Error("error marking task as done", slog.String("error", err.Error()))
		return
	}
	b.editTaskViewCallback(ctx, query, taskID)
}
