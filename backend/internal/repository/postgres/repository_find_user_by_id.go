package repository

import (
	"context"
	"proj/internal/entity"
)

func (s *PostgresStorage) FindByIdTg(ctx context.Context, userID int64) (int, error) {
	query := `SELECT user_id FROM user_telegram WHERE tg_id = $1`
	var id int
	err := s.DB.GetContext(ctx, &id, query, userID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PostgresStorage) FindByIdWeb(ctx context.Context, userID int) (*entity.UserWeb, error) {
	query := `SELECT u.username, up.email 
FROM users u 
JOIN user_passwords up ON u.id = up.user_id 
WHERE u.id = $1`
	user := entity.UserWeb{
		UserID: userID,
	}
	err := s.DB.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
