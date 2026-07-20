package usecase

import (
	"context"
	"proj/internal/entity"
)

type Task = entity.Task

type TaskRepository interface {
	Save(ctx context.Context, title string, userID int, done bool, description string) (id int, userTaskID int, err error)
	GetAllTasksByUserID(ctx context.Context, userID int) ([]Task, error)
	GetDeleteTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error)
	GetTask(ctx context.Context, taskID int) (*Task, error)
	GetTaskByUserID(ctx context.Context, userID int, user_task_id int) (int, error)
	Delete(ctx context.Context, id int) error
	DeleteForever(ctx context.Context, id int) error
	Recover(ctx context.Context, id int) error
	UpdateStatus(ctx context.Context, id int, status bool) error
	UpdateDescription(ctx context.Context, taskID int, newDesc string) error
}

type TaskUsecaseImpl struct {
	TaskRepo TaskRepository
}

func NewTaskUsecase(taskRepo TaskRepository) *TaskUsecaseImpl {
	return &TaskUsecaseImpl{TaskRepo: taskRepo}
}
