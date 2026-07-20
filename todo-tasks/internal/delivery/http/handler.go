package httpHandler

import (
	"context"
	"log/slog"
	"proj/internal/entity"
)

type TaskUsecase interface {
	CreateTask(ctx context.Context, title string, description string, userID int) (*entity.Task, error)
	DeleteTask(ctx context.Context, id int) error
	DeleteForeverTask(ctx context.Context, id int) error
	RecoverTask(ctx context.Context, id int) error
	MarkAsDone(ctx context.Context, id int, status bool) error
	UpdateDescription(ctx context.Context, taskID int, newDesc string) error
	GetTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error)
	GetTask(ctx context.Context, taskID int) (*entity.Task, error)
	GetTaskByUserID(ctx context.Context, userID int, user_task_id int) (int, error)
	GetRemovedTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error)
}

type TaskHandler struct {
	TaskUC TaskUsecase
	Logger *slog.Logger
}

func NewTaskHandler(taskUC TaskUsecase, logger *slog.Logger) *TaskHandler {
	return &TaskHandler{TaskUC: taskUC, Logger: logger}
}
