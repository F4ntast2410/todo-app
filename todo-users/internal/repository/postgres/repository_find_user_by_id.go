package repository

import (
	"context"
	"proj/internal/entity"
	pgmodel "proj/internal/repository/postgres/entityPostgres"
)

func (s *PostgresStorage) FindByIdTg(ctx context.Context, tgID int64) (int, error) {
	var userID int
	query := `SELECT user_id FROM user_telegram WHERE tg_id = $1`
	err := s.DB.GetContext(ctx, &userID, query, tgID)
	return userID, err
}

func (s *PostgresStorage) FindByIdWeb(ctx context.Context, userID int) (*entity.UserWeb, error) {
	query := `SELECT u.username, up.email, up.password_hash
FROM users u 
JOIN user_passwords up ON u.id = up.user_id 
WHERE u.id = $1`
	user := pgmodel.UserWeb{
		UserID: userID,
	}
	err := s.DB.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, err
	}

	return user.ToEntitiy(), nil
}
