package repository

import "context"

func (s *PostgresStorage) GetByUserID(ctx context.Context, userID int64) ([]Task, error) {
	var tasks []Task

	// Выбираем только те задачи, которые принадлежат конкретному ТГ-чату
	query := `SELECT id, title, done, user_id FROM tasks WHERE user_id = $1`

	err := s.DB.SelectContext(ctx, &tasks, query, userID)
	return tasks, err
}
