package repository

import (
	"context"
	"proj/internal/entity"
)

func (s *PostgresStorage) GetAllTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error) {
	var tasks []Task

	// Выбираем только те задачи, которые принадлежат конкретному ТГ-чату
	query := `SELECT id, title, done, user_id FROM tasks WHERE user_id = $1 AND deleted_at IS NULL ORDER BY user_task_id`

	err := s.DB.SelectContext(ctx, &tasks, query, userID)
	if err != nil {
		return nil, err
	}
	var entityTasks []entity.Task
	for _, task := range tasks {
		entityTasks = append(entityTasks, task.ToEntitiy())
	}
	return entityTasks, err
}

func (s *PostgresStorage) GetTask(ctx context.Context, taskID int) (*entity.Task, error) {
	var task Task

	query := `SELECT id, title, description, done, deleted_at FROM tasks WHERE id = $1`

	err := s.DB.GetContext(ctx, &task, query, taskID)
	if err != nil {
		return nil, err
	}
	entityTask := task.ToEntitiy()
	return &entityTask, err
}
