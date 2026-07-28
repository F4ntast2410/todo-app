package tgbot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleTaskList(ctx context.Context, userID int) (string, tgbotapi.InlineKeyboardMarkup) {
	tasks, err := b.taskUC.GetTasksByUserID(ctx, userID)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			b.logger.Error("database timeout during task creating")
			return "⚠️ База данных не ответила вовремя. Попробуйте еще раз.", tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Назад", "back_to_list"),
				),
			)
		}
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
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Добавить задачу📝", "create_task_step")))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Корзина🗑️", "trash_list")))
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
	var status string
	var status_text string
	var status_delete string
	var status_delete_query string
	button_switch_delete_status := tgbotapi.NewInlineKeyboardButtonData(status_text, fmt.Sprintf("mark_as_done:%d", taskID))
	switch task.Done {
	case true:
		status = "✅"
		status_text = "Отметить невыполненным❌️"
	case false:
		status = "❌️"
		status_text = "Отметить выполненным✅"
	}
	if task.DeletedAt == nil {
		button_switch_delete_status.Text = status_text
		status_delete = "Удалить задачу⛔"
		status_delete_query = "delete_task"
	} else {
		status_text = fmt.Sprintf("del_perm_confirm:%d", taskID)
		button_switch_delete_status.Text = "Удалить навсегда🗑️"
		button_switch_delete_status.CallbackData = &status_text
		status_delete = "Восстановить задачу🔄"
		status_delete_query = "recover_task"
	}

	text := fmt.Sprintf("Заголовок: %s\nСтатус выполнения: %s\nОписание: %s", task.Title, status, task.Description)
	keyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			button_switch_delete_status,
			tgbotapi.NewInlineKeyboardButtonData(status_delete, fmt.Sprintf("%s:%d", status_delete_query, taskID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Редактировать описание✏️", fmt.Sprintf("update_description:%d", taskID)),
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
