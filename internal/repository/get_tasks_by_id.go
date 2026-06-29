package repository

import (
	"context"
	"proj/internal/entity"
)

func (s *PostgresStorage) GetByUserID(ctx context.Context, userID int) ([]entity.Task, error) {
	var tasks []Task

	// Выбираем только те задачи, которые принадлежат конкретному ТГ-чату
	query := `SELECT id, title, done, user_id FROM tasks WHERE user_id = $1`

	err := s.DB.SelectContext(ctx, &tasks, query, userID)
	var entityTasks []entity.Task
	for _, task := range tasks {
		entityTasks = append(entityTasks, task.ToTask())
	}
	return entityTasks, err
}
