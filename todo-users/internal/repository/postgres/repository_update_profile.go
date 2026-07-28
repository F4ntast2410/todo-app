package repository

import (
	"context"
	"proj/internal/entity"
)

func (s *PostgresStorage) UpdateUsername(ctx context.Context, userID int, username string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE users SET username = $1 WHERE id = $2`, username, userID)
	return err
}
func (s *PostgresStorage) UpdateEmail(ctx context.Context, userID int, email entity.VerificationEmail) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE user_passwords SET email = $1 WHERE user_id = $2`, email, userID)
	return err
}

func (s *PostgresStorage) UpdatePasswordHash(ctx context.Context, userID int, passwordHash string) error {
	res, err := s.DB.ExecContext(ctx, `UPDATE user_passwords SET password_hash = $1 WHERE user_id = $2`, passwordHash, userID)
	if err != nil {
		return err
	}
	return checkAffected(res, "пользователь не найден")
}
