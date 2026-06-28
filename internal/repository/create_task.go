package repository

import "context"

func (s *PostgresStorage) Save(ctx context.Context, title string, userID int, done bool) (int, error) {
	query := `INSERT INTO tasks (user_id, user_task_id, title, done) 
VALUES (
    $1, 
    COALESCE((SELECT MAX(user_task_id) FROM tasks WHERE user_id = $1), 0) + 1, 
    $2,
	$3
) 
RETURNING user_task_id`
	var id int
	err := s.DB.GetContext(ctx, &id, query, userID, title, done)
	if err != nil {
		return -1, nil
	}
	return id, nil
}
