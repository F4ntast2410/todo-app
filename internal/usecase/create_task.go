package usecase

import (
	"context"
	"fmt"
)

func (uc *TaskUsecaseImpl) CreateTask(ctx context.Context, title string, userID int) (*Task, error) {
	if title == "" {
		return nil, fmt.Errorf("название не может быть пустым")
	}
	t := &Task{
		Title:  title,
		Done:   false,
		UserID: userID,
	}
	if err := uc.Repo.Save(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}
