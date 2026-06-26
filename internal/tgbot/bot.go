package tgbot

import (
	"context"
	"log/slog"

	"proj/internal/entity"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Определяем интерфейсы бизнес-логики, которые нужны боту
type TaskUsecase interface {
	CreateTask(ctx context.Context, title string, userID int) (*entity.Task, error)
	// Сюда потом допишем GetTasksByUserID, DeleteTask и т.д.
}
type UserUsecase interface {
	RegisterUserTg(ctx context.Context, user *entity.UserTg) error
	GetUserByTgID(ctx context.Context, user *entity.UserTg) error
}

type BotServer struct {
	bot    *tgbotapi.BotAPI
	taskUC TaskUsecase
	userUC UserUsecase
	logger *slog.Logger
}

func NewBotServer(token string, taskUC TaskUsecase, userUC UserUsecase, logger *slog.Logger) (*BotServer, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &BotServer{
		bot:    bot,
		taskUC: taskUC,
		userUC: userUC,
		logger: logger,
	}, nil
}
