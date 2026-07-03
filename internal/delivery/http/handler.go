package httpHandler

import (
	"context"
	"log/slog"
	"proj/internal/entity"
)

type TaskUsecase interface {
	CreateTask(ctx context.Context, title string, description string, userID int) (*entity.Task, error)
	DeleteTask(ctx context.Context, id int) error
	RecoverTask(ctx context.Context, id int) error
	MarkAsDone(ctx context.Context, id int, status bool) error
	UpdateDescription(ctx context.Context, taskID int, newDesc string) error
	GetTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error)
	GetTask(ctx context.Context, taskID int) (*entity.Task, error)
	GetRemovedTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error)
}

type TaskHandler struct {
	UC     TaskUsecase
	Logger *slog.Logger
}
