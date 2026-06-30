package repository

import "context"

func (s *PostgresStorage) UpdateDescription(ctx context.Context, taskID int, newDesc string) error {
	query := `UPDATE tasks SET description = $1 WHERE id = $2`

	// Используем ExecContext, так как нам не нужно сканировать возвращаемые строки
	_, err := s.DB.ExecContext(ctx, query, newDesc, taskID)
	return err
}
