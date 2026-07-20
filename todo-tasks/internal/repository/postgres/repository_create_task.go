package repository

import "context"

func (s *PostgresStorage) Save(ctx context.Context, title string, userID int, done bool, description string) (id int, userTaskID int, err error) {
	query := `INSERT INTO tasks (user_id, user_task_id, title, done, description) 
VALUES (
    $1, 
    COALESCE((SELECT MAX(user_task_id) FROM tasks WHERE user_id = $1), 0) + 1, 
    $2, $3, $4
) 
RETURNING id, user_task_id`

	err = s.DB.QueryRowContext(ctx, query, userID, title, done, description).Scan(&id, &userTaskID)
	if err != nil {
		return 0, 0, err
	}
	return id, userTaskID, nil
}
