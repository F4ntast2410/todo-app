package repository

import (
	"context"
	"proj/internal/entity"
	pgmodel "proj/internal/repository/postgres/entityPostgres"
)

func (s *PostgresStorage) FindByEmail(ctx context.Context, email entity.VerificationEmail) (*entity.UserWeb, error) {
	query := `SELECT password_hash, user_id, email FROM user_passwords WHERE email = $1`
	var user pgmodel.UserWeb
	err := s.DB.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return user.ToEntitiy(), nil
}
