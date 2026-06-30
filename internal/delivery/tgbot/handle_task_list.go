package tgbot

import (
	"context"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleTaskList(ctx context.Context, userID int) (string, tgbotapi.InlineKeyboardMarkup) {
	tasks, err := b.taskUC.GetTasksByUserID(ctx, userID)
	if err != nil {
		b.logger.Error("error getting task", slog.String("error", err.Error()))
		return "⚠️ Задача не найдена", tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", "back_to_list"),
			),
		)
	}
	text := "Ваши задачи:"
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, t := range tasks {
		done := "❌️"
		if t.Done {
			done = "✅"
		}
		btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %s\n", t.Title, done), fmt.Sprintf("task_view:%d", t.ID))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Добавить задачу", "create_task_step")))
	return text, tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (b *BotServer) sendTaskList(ctx context.Context, msg *tgbotapi.Message) {
	user, err := b.userUC.GetUserByTgID(ctx, msg.From.ID, msg.From.UserName)
	if err != nil {
		b.logger.Error("error getting user id", slog.String("error", err.Error()))
		return
	}
	text, kb := b.handleTaskList(ctx, user.UserID)
	message := tgbotapi.NewMessage(msg.From.ID, text)
	message.ReplyMarkup = kb
	b.Send(message)
}

func (b *BotServer) editTaskList(ctx context.Context, query *tgbotapi.CallbackQuery) {
	user, err := b.userUC.GetUserByTgID(ctx, query.From.ID, query.Message.From.UserName)
	if err != nil {
		b.logger.Error("error getting user id", slog.String("error", err.Error()))
		return
	}
	text, kb := b.handleTaskList(ctx, user.UserID)
	edit := tgbotapi.NewEditMessageTextAndMarkup(query.From.ID, query.Message.MessageID, text, kb)
	b.Send(edit)
}

func (b *BotServer) handleTaskView(ctx context.Context, taskID int) (string, tgbotapi.InlineKeyboardMarkup) {
	task, err := b.taskUC.GetTask(ctx, taskID)
	if err != nil {
		b.logger.Error("error getting task", slog.String("error", err.Error()))
		return "⚠️ Задача не найдена", tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", "back_to_list"),
			),
		)
	}
	var done string
	var rev_done_text string
	var isdeleted string
	var delQuery string
	switch task.Done {
	case true:
		done = "✅"
		rev_done_text = "невыполненным"
	case false:
		done = "❌️"
		rev_done_text = "выполненным"
	}
	if task.DeletedAt == nil {
		isdeleted = "Удалить задачу"
		delQuery = "delete_task"
	} else {
		isdeleted = "Восстановить задачу"
		delQuery = "recover_task"
	}
	text := fmt.Sprintf("Заголовок: %s\nСтатус выполнения: %s\nОписание: %s", task.Title, done, task.Description)
	keyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Отметить %s", rev_done_text), fmt.Sprintf("mark_as_done:%d", taskID)),
			tgbotapi.NewInlineKeyboardButtonData(isdeleted, fmt.Sprintf("%s:%d", delQuery, taskID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Редактировать описание", fmt.Sprintf("update_description:%d", taskID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back_to_list"),
		),
	)
	return text, keyboardMarkup

}

func (b *BotServer) sendTaskViewMessage(ctx context.Context, msg *tgbotapi.Message, taskID int) {
	text, kb := b.handleTaskView(ctx, taskID)
	message := tgbotapi.NewMessage(msg.From.ID, text)
	message.ReplyMarkup = kb
	b.Send(message)
}

func (b *BotServer) editTaskViewCallback(ctx context.Context, query *tgbotapi.CallbackQuery, taskID int) {
	text, kb := b.handleTaskView(ctx, taskID)
	edit := tgbotapi.NewEditMessageTextAndMarkup(query.From.ID, query.Message.MessageID, text, kb)
	b.Send(edit)
}
