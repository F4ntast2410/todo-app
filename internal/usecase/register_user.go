package usecase

import (
	"context"
	"errors"
	"proj/internal/entity"
)

var ErrUserAlreadyExists = errors.New("user already exists")

func (uc *UserUsecaseImpl) RegisterUserWeb(ctx context.Context, user *entity.UserWeb) error {
	exists, err := uc.Repo.ExistsWeb(ctx, user)
	if err != nil {
		return err
	}
	if exists == false {
		return uc.Repo.CreateUserWeb(ctx, user)
	} else {
		return ErrUserAlreadyExists
	}
}

func (uc *UserUsecaseImpl) RegisterUserTg(ctx context.Context, user *entity.UserTg) error {
	exists, err := uc.Repo.ExistsTg(ctx, user)
	if err != nil {
		return err
	}
	if exists == false {
		return uc.Repo.CreateUserTg(ctx, user)
	} else {
		return ErrUserAlreadyExists
	}
}
