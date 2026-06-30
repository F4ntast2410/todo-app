package tgbot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handleTaskCreation(ctx context.Context, msg *tgbotapi.Message) {
	message := tgbotapi.NewMessage(msg.Chat.ID, "")

	taskDescription := msg.Text

	b.sessionCache.Mu.Lock()
	taskTitle := b.sessionCache.Cache[msg.From.ID].TaskTitle
	session := b.sessionCache.Cache[msg.From.ID]
	session.TaskTitle = ""
	session.State = StateIdle
	b.sessionCache.Cache[msg.From.ID] = session
	b.sessionCache.Mu.Unlock()

	// Передаем этот контекст дальше по всем слоям до самой БД!
	user, err := b.userUC.GetUserByTgID(ctx, msg.From.ID, msg.From.UserName)
	if err != nil {
		err = b.userUC.RegisterUserTg(ctx, msg.From.ID, msg.From.UserName)
		if err != nil {
			b.logger.Error("failed to register user", slog.String("error", err.Error()))
			message.Text = "⚠️ Не удалось завершить регистрацию."
			b.Send(message)
			return
		}
		message.Text = "Привет! Ты успешно зарегистрирован в система, повтори свой запрос"
		b.Send(message)
		return
	}
	task, err := b.taskUC.CreateTask(ctx, taskTitle, taskDescription, user.UserID)
	if err != nil {
		// Если контекст завершился по таймауту
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			b.logger.Error("database timeout during task creating")
			message.Text = "⚠️ База данных не ответила вовремя. Попробуйте еще раз."
			b.Send(message)
			return
		}
		message.Text = fmt.Sprintf("⚠️ Ошибка: %s", err.Error())
		b.Send(message)
		b.logger.Error("failed to create task", slog.String("error", err.Error()))
		return
	}
	b.sendTaskViewMessage(ctx, msg, task.ID)
}

func (b *BotServer) handleGetTitleMessage(msg *tgbotapi.Message) {
	message := tgbotapi.NewMessage(msg.From.ID, "Опишите задачу: ")
	maxLen := 40
	if len(msg.Text) > maxLen {
		message.Text = fmt.Sprintf("Размер названия задачи не должен превышать %d символов", maxLen)
		b.Send(message)
		b.sessionCache.Mu.Lock()
		session := b.sessionCache.Cache[msg.From.ID]
		session.State = StateIdle
		b.sessionCache.Cache[msg.From.ID] = session
		b.sessionCache.Mu.Unlock()
		return
	}
	b.sessionCache.Mu.Lock()
	session := b.sessionCache.Cache[msg.From.ID]
	session.TaskTitle = msg.Text
	session.State = StateWaitingTaskDescription
	b.sessionCache.Cache[msg.From.ID] = session
	b.sessionCache.Mu.Unlock()
	b.Send(message)

}

func (b *BotServer) handleGetTitleCallback(query *tgbotapi.CallbackQuery) {
	b.sessionCache.Mu.Lock()
	session := b.sessionCache.Cache[query.From.ID]
	session.State = StateWaitingTaskTitle
	b.sessionCache.Cache[query.From.ID] = session
	b.sessionCache.Mu.Unlock()
	message := tgbotapi.NewMessage(query.From.ID, "Введите название задачи: ")
	b.Send(message)
}
