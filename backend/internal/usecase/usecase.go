package usecase

import (
	"context"
	"proj/internal/entity"
	"time"
)

type Task = entity.Task

type TaskRepository interface {
	Save(ctx context.Context, title string, userID int, done bool, description string) (int, error)
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

type UserUsecaseImpl struct {
	UserRepo UserRepository
	AuthRepo VerificationRepository
}

type UserRepository interface {
	CreateUserWeb(ctx context.Context, email string, passwordHash string, username string) error
	CreateUserTg(ctx context.Context, ID int64, username string) error
	ExistsWeb(ctx context.Context, email string) (bool, error)
	ExistsTg(ctx context.Context, ID int64) (bool, error)
	FindByIdTg(ctx context.Context, userID int64) (int, error)
	FindByIdWeb(ctx context.Context, userID int) (*entity.UserWeb, error)
	FindByEmail(ctx context.Context, email entity.VerificationEmail) (*entity.UserWeb, error)
	CreateSession(ctx context.Context, userID int) (string, error)
	FindUserIDBySession(ctx context.Context, sessionID string) (int, error)
}

type VerificationRepository interface {
	Save2FA(ctx context.Context, email entity.VerificationEmail, unique_code string, time_duration time.Duration) error
	Get2FA(ctx context.Context, email entity.VerificationEmail) (string, error)
	Delete2FA(ctx context.Context, email entity.VerificationEmail) error
}

func NewUserUsecase(userRepo UserRepository, authRepo VerificationRepository) *UserUsecaseImpl {
	return &UserUsecaseImpl{UserRepo: userRepo, AuthRepo: authRepo}
}
func NewTaskUsecase(taskRepo TaskRepository) *TaskUsecaseImpl {
	return &TaskUsecaseImpl{TaskRepo: taskRepo}
}
