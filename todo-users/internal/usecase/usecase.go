package usecase

import (
	"context"
	"proj/internal/entity"
	"time"
)

type UserUsecaseImpl struct {
	UserRepo UserRepository
	AuthRepo VerificationRepository
	Mailer   Mailer
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
	UpdateProfile(ctx context.Context, userID int, username, email string) error
	UpdatePasswordHash(ctx context.Context, userID int, passwordHash string) error
	LinkTelegramAccount(ctx context.Context, userID int, tgID int64, username string) error
	GetTelegramLinkByUserID(ctx context.Context, userID int) (*entity.TelegramLink, error)
}

type VerificationRepository interface {
	Save2FA(ctx context.Context, email entity.VerificationEmail, unique_code string, time_duration time.Duration) error
	Get2FA(ctx context.Context, email entity.VerificationEmail) (string, error)
	Delete2FA(ctx context.Context, email entity.VerificationEmail) error
}
type Mailer interface {
	Send2FACode(to string, code string) error
}

func NewUserUsecase(userRepo UserRepository, authRepo VerificationRepository, mailer Mailer) *UserUsecaseImpl {
	return &UserUsecaseImpl{UserRepo: userRepo, AuthRepo: authRepo, Mailer: mailer}
}
