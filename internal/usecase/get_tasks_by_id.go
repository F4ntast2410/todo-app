package usecase

import (
	"context"
)

func (uc *TaskUsecaseImpl) GetTasksByUserID(ctx context.Context, userID int) ([]Task, error) {
	tasks, err := uc.TaskRepo.GetAllTasksByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil

}

func (uc *TaskUsecaseImpl) GetTask(ctx context.Context, taskID int) (*Task, error) {
	task, err := uc.TaskRepo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return task, nil

}
