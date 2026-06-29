package httpHandler

import (
	"context"
	"log/slog"
	"proj/internal/entity"
)

type TaskUsecase interface {
	CreateTask(ctx context.Context, title string, description string, userID int) (*entity.Task, error)
}

type TaskHandler struct {
	UC     TaskUsecase
	Logger *slog.Logger
}
