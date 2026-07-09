package repository

import (
	"context"
	"time"
)

func (s *PostgresStorage) CreateSession(ctx context.Context, userID int) (string, error) {
	query := `INSERT INTO user_sessions (user_id, expires_at) VALUES ($1, $2) RETURNING id`
	var id string
	err := s.DB.GetContext(ctx, &id, query, userID, time.Now().Add(24*time.Hour))
	if err != nil {
		return "", err
	}
	return id, nil
}
