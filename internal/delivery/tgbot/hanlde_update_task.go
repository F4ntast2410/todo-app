package tgbot

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotServer) handlerUpdateDecriptionCallback(query *tgbotapi.CallbackQuery, taskID int) {
	b.sessionCache.Mu.Lock()
	session := b.sessionCache.Cache[query.From.ID]
	session.State = StateWaitingNewTaskDescription
	session.TaskID = taskID
	b.sessionCache.Cache[query.From.ID] = session
	b.sessionCache.Mu.Unlock()
	message := tgbotapi.NewMessage(query.From.ID, "Введите новое описание задачи: ")
	b.Send(message)
}

func (b *BotServer) handlerUpdateDecription(ctx context.Context, msg *tgbotapi.Message) {
	b.sessionCache.Mu.RLock()
	taskID := b.sessionCache.Cache[msg.From.ID].TaskID
	session := b.sessionCache.Cache[msg.From.ID]
	session.State = StateIdle
	session.TaskID = 0
	b.sessionCache.Cache[msg.From.ID] = session
	b.sessionCache.Mu.RUnlock()
	newDesc := msg.Text
	err := b.taskUC.UpdateDescription(ctx, taskID, newDesc)
	if err != nil {
		b.logger.Error("error updating task", slog.String("error: ", err.Error()))
		return
	}

	b.sendTaskViewMessage(ctx, msg, taskID)
}
