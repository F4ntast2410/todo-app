package httpHandler

import (
	"context"
	"log/slog"
	"proj/internal/entity"
)

type UserUsecase interface {
	RegisterUserWeb(ctx context.Context, email, username, password string) error
	LoginUserWeb(ctx context.Context, email entity.VerificationEmail, password string) (*entity.LoginResult, error)
	Verify2FA(ctx context.Context, email entity.VerificationEmail, inputCode string) (*entity.LoginResult, error)
	GetUserIDBySession(ctx context.Context, sessionID string) (int, error)
	GetUserByWebID(ctx context.Context, userID int) (*entity.UserWeb, error)
	UpdateProfile(ctx context.Context, userID int, username, email string) error
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error
	GetTelegramLink(ctx context.Context, userID int) (*entity.TelegramLink, error)
	LinkTelegram(ctx context.Context, userID int, tgID int64, tgUsername string) error
	LoginByTelegram(ctx context.Context, tgID int64, tgUsername string) (*entity.LoginResult, error)
}

type UserHandler struct {
	UserUC   UserUsecase
	Logger   *slog.Logger
	BotToken string
}

func NewUserHandler(userUC UserUsecase, logger *slog.Logger, botToken string) *UserHandler {
	return &UserHandler{UserUC: userUC, Logger: logger, BotToken: botToken}
}
