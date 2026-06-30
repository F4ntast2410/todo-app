package usecase

import (
	"context"
)

func (uc *TaskUsecaseImpl) DeleteTask(ctx context.Context, id int) error {
	return uc.TaskRepo.Delete(ctx, id)
}

func (uc *TaskUsecaseImpl) RecoverTask(ctx context.Context, id int) error {
	return uc.TaskRepo.Recover(ctx, id)
}
