package tgbot

import (
	"context"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleTrashList(ctx context.Context, query *tgbotapi.CallbackQuery) {
	user, err := b.userUC.GetUserByTgID(ctx, query.From.ID, query.From.UserName)
	if err != nil {
		b.logger.Error("error getting user id", slog.String("error", err.Error()))
		return
	}
	tasks, err := b.taskUC.GetRemovedTasksByUserID(ctx, user.UserID)
	if err != nil {
		b.logger.Error("error getting task", slog.String("error", err.Error()))
		return
	}
	text := "Ваши задачи:"
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, t := range tasks {
		var status string
		switch t.Done {
		case true:
			status = "✅"
		case false:
			status = "❌️"
		}
		btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %s\n", t.Title, status), fmt.Sprintf("task_view:%d", t.ID))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад", "back_to_list")))
	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	edit := tgbotapi.NewEditMessageTextAndMarkup(query.From.ID, query.Message.MessageID, text, kb)
	b.Send(edit)
}
