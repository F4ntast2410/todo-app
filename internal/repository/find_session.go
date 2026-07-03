package repository

import (
	"context"
	"time"
)

func (s *PostgresStorage) FindUserIDBySession(ctx context.Context, sessionID string) (int, error) {
	query := `SELECT user_id FROM user_sessions WHERE id = $1 AND expires_at > $2`
	var userID int
	// Проверяем, что сессия существует и время её жизни больше, чем NOW()
	err := s.DB.GetContext(ctx, &userID, query, sessionID, time.Now())
	if err != nil {
		return 0, err
	}
	return userID, nil
}
