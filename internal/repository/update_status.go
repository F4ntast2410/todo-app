package repository

import "context"

func (s *PostgresStorage) UpdateStatus(ctx context.Context, id int, done bool) error {
	query := `UPDATE tasks SET done = $1 WHERE id = $2`

	// Используем ExecContext, так как нам не нужно сканировать возвращаемые строки
	_, err := s.DB.ExecContext(ctx, query, done, id)
	return err
}
