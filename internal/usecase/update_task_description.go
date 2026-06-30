package usecase

import "context"

func (uc TaskUsecaseImpl) UpdateDescription(ctx context.Context, taskID int, newDesc string) error {
	return uc.TaskRepo.UpdateDescription(ctx, taskID, newDesc)
}
