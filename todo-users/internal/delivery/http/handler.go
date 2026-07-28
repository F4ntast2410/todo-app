package httpHandler

import (
	"context"
	"log/slog"
	"proj/internal/entity"
)

type UserUsecase interface {
	RegisterUserWeb(ctx context.Context, email entity.VerificationEmail, username, password, inputCode string) (*entity.AuthResult, error)
	RegisterUserWebSend2FA(ctx context.Context, email entity.VerificationEmail) (*entity.AuthResult, error)
	LoginUserWeb(ctx context.Context, email entity.VerificationEmail, password string) (*entity.AuthResult, error)
	Verify2FA(ctx context.Context, email entity.VerificationEmail, inputCode string) (*entity.AuthResult, error)
	GetUserIDBySession(ctx context.Context, sessionID string) (int, error)
	GetUserByWebID(ctx context.Context, userID int) (*entity.UserWeb, error)
	UpdateUsername(ctx context.Context, userID int, username string) error
	UpdateEmail(ctx context.Context, newEmail, email entity.VerificationEmail) error
	UpdateEmailVerify(ctx context.Context, userID int, newEmail, email entity.VerificationEmail, inputCode string) error
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error
	GetTelegramLink(ctx context.Context, userID int) (*entity.TelegramLink, error)
	LinkTelegram(ctx context.Context, userID int, tgID int64, tgUsername string) error
	LoginByTelegram(ctx context.Context, tgID int64, tgUsername string) (*entity.AuthResult, error)
	SetEmail(ctx context.Context, email entity.VerificationEmail) error
	SetEmailPassword(ctx context.Context, email entity.VerificationEmail, password string, userID int) error
	Verify2FACode(ctx context.Context, email entity.VerificationEmail, inputCode string) error
}

type UserHandler struct {
	UserUC   UserUsecase
	Logger   *slog.Logger
	BotToken string
}

func NewUserHandler(userUC UserUsecase, logger *slog.Logger, botToken string) *UserHandler {
	return &UserHandler{UserUC: userUC, Logger: logger, BotToken: botToken}
}
