package usecase

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (uc *UserUsecaseImpl) UpdateProfile(ctx context.Context, userID int, username, email string) error {
	return uc.UserRepo.UpdateProfile(ctx, userID, username, email)
}

func (uc *UserUsecaseImpl) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	user, err := uc.UserRepo.FindByIdWeb(ctx, userID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("неверный текущий пароль")
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return uc.UserRepo.UpdatePasswordHash(ctx, userID, string(newHash))
}
