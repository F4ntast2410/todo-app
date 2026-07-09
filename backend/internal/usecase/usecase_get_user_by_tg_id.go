package usecase

import (
	"context"
	"proj/internal/entity"
)

func (uc *UserUsecaseImpl) GetUserByTgID(ctx context.Context, userID int64, username string) (*entity.UserTg, error) {
	id, err := uc.UserRepo.FindByIdTg(ctx, userID)
	if err != nil {
		return nil, err
	}
	user := entity.UserTg{
		ID:       userID,
		UserID:   id,
		Username: username}

	return &user, nil
}

func (uc *UserUsecaseImpl) GetUserByWebID(ctx context.Context, userID int) (*entity.UserWeb, error) {
	user, err := uc.UserRepo.FindByIdWeb(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
