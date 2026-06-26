package usecase

import (
	"context"
	"proj/internal/entity"
	"proj/internal/repository"
)

type Task = repository.Task

type TaskRepository interface {
	Save(ctx context.Context, task *Task) error
	// GetByUserID(ctx context.Context, userID int64) ([]Task, error)
}

type TaskUsecaseImpl struct {
	Repo TaskRepository
}

type UserRepository interface {
	CreateUserWeb(ctx context.Context, user *entity.UserWeb) error
	CreateUserTg(ctx context.Context, user *entity.UserTg) error
	ExistsWeb(ctx context.Context, user *entity.UserWeb) (bool, error)
	ExistsTg(ctx context.Context, user *entity.UserTg) (bool, error)
	FindByIdTg(ctx context.Context, user *entity.UserTg) error
}

type UserUsecaseImpl struct {
	Repo UserRepository
}
