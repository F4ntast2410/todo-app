package repository

import (
	"context"
)

func (s *PostgresStorage) Delete(ctx context.Context, id int) error {
	query := `UPDATE tasks SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := s.DB.ExecContext(ctx, query, id)
	return err
}

func (s *PostgresStorage) Recover(ctx context.Context, id int) error {
	query := `UPDATE tasks SET deleted_at = NULL WHERE id = $1 AND deleted_at IS NOT NULL`
	_, err := s.DB.ExecContext(ctx, query, id)
	return err
}
