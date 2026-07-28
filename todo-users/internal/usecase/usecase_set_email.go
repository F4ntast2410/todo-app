package usecase

import (
	"context"
	"proj/internal/entity"
	customErrors "proj/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

func (uc *UserUsecaseImpl) SetEmail(ctx context.Context, email entity.VerificationEmail) error {
	exists, err := uc.UserRepo.ExistsWeb(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return customErrors.ErrUserAlreadyExists
	}
	return uc.Send2FACode(ctx, email)
}
func (uc *UserUsecaseImpl) SetEmailPassword(ctx context.Context, email entity.VerificationEmail, password string, userID int) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = uc.UserRepo.LinkWebAccount(ctx, userID, email, string(passwordHash))
	return err
}
