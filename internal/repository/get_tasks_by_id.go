package repository

import (
	"context"
	"proj/internal/entity"
)

func (s *PostgresStorage) GetAllTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error) {
	var tasks []Task

	// Выбираем только те задачи, которые принадлежат конкретному ТГ-чату
	query := `SELECT id, title, done, user_id, description, user_task_id FROM tasks WHERE user_id = $1 AND deleted_at IS NULL ORDER BY user_task_id`

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

func (s *PostgresStorage) GetDeleteTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error) {
	var tasks []Task

	// Выбираем только те задачи, которые принадлежат конкретному ТГ-чату
	query := `SELECT id, title, done, user_id, description, deleted_at, user_task_id FROM tasks WHERE user_id = $1 AND deleted_at IS NOT NULL ORDER BY user_task_id`

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

func (s *PostgresStorage) GetTaskByUserID(ctx context.Context, userID int, user_task_id int) (int, error) {

	query := `SELECT id FROM tasks WHERE user_id = $1 AND user_task_id = $2`
	var id int
	err := s.DB.GetContext(ctx, &id, query, userID, user_task_id)
	if err != nil {
		return 0, err
	}
	return id, err
}
