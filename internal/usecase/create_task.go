package usecase

import (
	"context"
	"fmt"
	"proj/internal/usecase/entity"
)

func (uc *UsecaseImpl) CreateTask(ctx context.Context, title string, userID int) (*entity.Task, error) {
	if title == "" {
		return nil, fmt.Errorf("название не может быть пустым")
	}
	t := &entity.Task{
		Title:  title,
		Done:   false,
		UserID: userID,
	}
	id, err := uc.TaskRepo.Save(ctx, t.Title, t.UserID, t.Done)
	if err != nil {
		return nil, err
	}
	t.ID = id
	return t, nil
}
