package repository

import (
	"context"
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
