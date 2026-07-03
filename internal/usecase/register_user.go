package usecase

import (
	"context"
	customErrors "proj/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

func (uc *UserUsecaseImpl) RegisterUserWeb(ctx context.Context, email, username, password string) error {
	exists, err := uc.UserRepo.ExistsWeb(ctx, email)
	if err != nil {
		return err
	}
	if !exists {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		err = uc.UserRepo.CreateUserWeb(ctx, email, string(passwordHash), username)
		if err != nil {
			return err
		}
		return nil
	} else {
		return customErrors.ErrUserAlreadyExists
	}
}

func (uc *UserUsecaseImpl) RegisterUserTg(ctx context.Context, ID int64, username string) error {
	exists, err := uc.UserRepo.ExistsTg(ctx, ID)
	if err != nil {
		return err
	}
	if !exists {
		err := uc.UserRepo.CreateUserTg(ctx, ID, username)
		if err != nil {
			return err
		}
		return nil
	} else {
		return customErrors.ErrUserAlreadyExists
	}
}
