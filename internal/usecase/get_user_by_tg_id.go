package usecase

import (
	"context"
	"proj/internal/entity"
)

func (uc *UserUsecaseImpl) GetUserByTgID(ctx context.Context, user *entity.UserTg) error {
	return uc.Repo.FindByIdTg(ctx, user)
}
