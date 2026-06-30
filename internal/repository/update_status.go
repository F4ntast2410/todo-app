package repository

import "context"

func (s *PostgresStorage) UpdateStatus(ctx context.Context, id int, status bool) error {
	query := `UPDATE tasks SET done = $1 WHERE id = $2`

	// Используем ExecContext, так как нам не нужно сканировать возвращаемые строки
	result, err := s.DB.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}
	return checkAffected(result, "task not found")
}
