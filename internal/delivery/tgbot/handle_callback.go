package tgbot

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleCallback(query *tgbotapi.CallbackQuery) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := strings.Split(query.Data, ":")

	switch data[0] {
	case "create_task_step":
		b.handleGetTitleCallback(query)
	case "task_view":
		taskID, err := strconv.Atoi(data[1])
		if err != nil {
			b.logger.Error("error converting string to int", slog.String("error", err.Error()))
			return
		}
		b.editTaskViewCallback(ctx, query, taskID)
	case "back_to_list":
		b.editTaskList(ctx, query)
	case "mark_as_done":
		taskID, err := strconv.Atoi(data[1])
		if err != nil {
			b.logger.Error("error converting string to int", slog.String("error", err.Error()))
			return
		}
		b.handlerMarkAsDone(ctx, query, taskID)
	case "delete_task":
		taskID, err := strconv.Atoi(data[1])
		if err != nil {
			b.logger.Error("error converting string to int", slog.String("error", err.Error()))
			return
		}
		b.handlerDeleteTask(ctx, query, taskID)
	case "recover_task":
		taskID, err := strconv.Atoi(data[1])
		if err != nil {
			b.logger.Error("error converting string to int", slog.String("error", err.Error()))
			return
		}
		b.handlerRecoverTask(ctx, query, taskID)
	case "update_description":
		taskID, err := strconv.Atoi(data[1])
		if err != nil {
			b.logger.Error("error converting string to int", slog.String("error", err.Error()))
			return
		}
		b.handlerUpdateDecriptionCallback(query, taskID)
	}
}
