package repository

import (
	"context"
	"time"
)

func (s *PostgresStorage) CreateSession(ctx context.Context, userID int, tokenHash string) error {
	query := `INSERT INTO user_sessions (user_id, expires_at, token_hash) VALUES ($1, $2, $3) RETURNING id`
	_, err := s.DB.ExecContext(ctx, query, userID, time.Now().Add(24*time.Hour), tokenHash)
	return err
}
