package usecase

import (
	"context"
	"proj/internal/entity"
)

type Task = entity.Task

type TaskRepository interface {
	Save(ctx context.Context, title string, userID int, done bool, description string) (int, error)
	GetByUserID(ctx context.Context, userID int) ([]Task, error)
}

type TaskUsecaseImpl struct {
	TaskRepo TaskRepository
}

type UserUsecaseImpl struct {
	UserRepo UserRepository
}

type UserRepository interface {
	CreateUserWeb(ctx context.Context, email string, passwordHash string, username string) error
	CreateUserTg(ctx context.Context, ID int64, username string) error
	ExistsWeb(ctx context.Context, email string) (bool, error)
	ExistsTg(ctx context.Context, ID int64) (bool, error)
	FindByIdTg(ctx context.Context, userID int64) (int, error)
}
