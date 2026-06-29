package usecase

import (
	"context"
	"fmt"
	"proj/internal/entity"
)

func (uc *TaskUsecaseImpl) CreateTask(ctx context.Context, title string, description string, userID int) (*entity.Task, error) {
	if title == "" {
		return nil, fmt.Errorf("название не может быть пустым")
	}
	t := &entity.Task{
		Title:       title,
		Done:        false,
		UserID:      userID,
		Description: description,
	}
	id, err := uc.TaskRepo.Save(ctx, t.Title, t.UserID, t.Done, t.Description)
	if err != nil {
		return nil, err
	}
	t.ID = id
	return t, nil
}
