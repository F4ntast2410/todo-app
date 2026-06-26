package repository

import "context"

func (s *PostgresStorage) Save(ctx context.Context, task *Task) error {
	query := `INSERT INTO tasks (user_id, user_task_id, title, done) 
VALUES (
    $1, 
    COALESCE((SELECT MAX(user_task_id) FROM tasks WHERE user_id = $1), 0) + 1, 
    $2,
	$3
) 
RETURNING user_task_id`
	return s.DB.GetContext(ctx, &task.ID, query, task.UserID, task.Title, task.Done)
}
