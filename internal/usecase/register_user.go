package usecase

import (
	"context"
	"errors"
)

var ErrUserAlreadyExists error = errors.New("user already exists")

func (uc *UsecaseImpl) RegisterUserWeb(ctx context.Context, email string, username string, passwordHash string) error {
	exists, err := uc.UserRepo.ExistsWeb(ctx, email)
	if err != nil {
		return err
	}
	if exists == false {
		err := uc.UserRepo.CreateUserWeb(ctx, email, username, passwordHash)
		if err != nil {
			return err
		}
		return nil
	} else {
		return ErrUserAlreadyExists
	}
}

func (uc *UsecaseImpl) RegisterUserTg(ctx context.Context, ID int64, username string) error {
	exists, err := uc.UserRepo.ExistsTg(ctx, ID)
	if err != nil {
		return err
	}
	if exists == false {
		err := uc.UserRepo.CreateUserTg(ctx, ID, username)
		if err != nil {
			return err
		}
		return nil
	} else {
		return ErrUserAlreadyExists
	}
}
