package repository

import (
	"context"
	"proj/internal/entity"
)

func (s *PostgresStorage) ExistsWeb(ctx context.Context, user *entity.UserWeb) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM user_password WHERE email = $1)`

	// Используем GetContext, так как запрос ВСЕГДА возвращает ровно одну строку с true/false
	err := s.DB.GetContext(ctx, &exists, query, user.Email)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *PostgresStorage) ExistsTg(ctx context.Context, user *entity.UserTg) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM user_telegram WHERE tg_id = $1)`

	// Используем GetContext, так как запрос ВСЕГДА возвращает ровно одну строку с true/false
	err := s.DB.GetContext(ctx, &exists, query, user.ID)
	if err != nil {
		return false, err
	}

	return exists, nil
}
