package usecase

import (
	"context"
)

func (uc *TaskUsecaseImpl) MarkAsDone(ctx context.Context, id int, status bool) error {
	return uc.TaskRepo.UpdateStatus(ctx, id, status)
}
