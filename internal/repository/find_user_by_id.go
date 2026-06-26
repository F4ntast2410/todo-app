package repository

import (
	"context"
	"proj/internal/entity"
)

func (s *PostgresStorage) FindByIdTg(ctx context.Context, user *entity.UserTg) error {
	query := `SELECT user_id FROM user_telegram WHERE tg_id = $1`
	err := s.DB.GetContext(ctx, &user.UserID, query, user.ID)
	if err != nil {
		return err
	}
	return nil
}
