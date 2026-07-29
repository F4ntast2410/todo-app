package usecase

import (
	"context"
	"errors"
	"proj/internal/entity"
	customErrors "proj/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

func (uc *UserUsecaseImpl) UpdateUsername(ctx context.Context, userID int, username string) error {
	return uc.UserRepo.UpdateUsername(ctx, userID, username)
}
func (uc *UserUsecaseImpl) UpdateEmail(ctx context.Context, newEmail, email entity.VerificationEmail) error {
	exists, err := uc.UserRepo.ExistsWeb(ctx, newEmail)
	if err != nil {
		return err
	}
	if exists {
		return customErrors.ErrUserAlreadyExists
	} else {
		err = uc.Send2FACode(ctx, email)
		if err != nil {
			return err
		}
	}
	return nil
}

func (uc *UserUsecaseImpl) UpdateEmailVerify(ctx context.Context, userID int, newEmail, email entity.VerificationEmail, inputCode string) error {
	err := uc.Verify2FACode(ctx, email, inputCode)
	if err != nil {
		return err
	}
	err = uc.UserRepo.UpdateEmail(ctx, userID, newEmail)
	return err
}

func (uc *UserUsecaseImpl) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	user, err := uc.UserRepo.FindByIdWeb(ctx, userID)
	if err != nil {
		return err
	}
	if !user.HasPassword {
		return errors.New("user not exists")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return customErrors.ErrInccorectPassword
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return uc.UserRepo.UpdatePasswordHash(ctx, userID, string(newHash))
}
