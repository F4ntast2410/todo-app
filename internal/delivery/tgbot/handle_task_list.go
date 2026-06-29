package tgbot

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleTaskList(ctx context.Context, msg *tgbotapi.Message) {
	user, err := b.userUC.GetUserByTgID(ctx, msg.From.ID, msg.From.UserName)
	if err != nil {
		b.logger.Error("error getting user by id", slog.String("error:", err.Error()))
		return
	}
	tasks, err := b.taskUC.GetTasksByUserID(ctx, user.UserID)
	if err != nil {
		b.logger.Error("error getting tasks by user id", slog.String("error:", err.Error()))
		return
	}
	str := "Вот твой список задач:\n"
	for _, task := range tasks {
		if task.Done {
			str += "• " + task.Title + " ✔\n"
			continue
		}
		str += "• " + task.Title + " X\n"
	}
	query := "create_task_step"
	keyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить задачу", query),
		),
	)
	message := tgbotapi.NewMessage(msg.Chat.ID, str)
	message.ReplyMarkup = &keyboardMarkup
	b.Send(message)
}
