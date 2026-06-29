package usecase

import (
	"context"
	"proj/internal/entity"
)

func (uc *TaskUsecaseImpl) GetTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error) {
	tasks, err := uc.TaskRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil

}
