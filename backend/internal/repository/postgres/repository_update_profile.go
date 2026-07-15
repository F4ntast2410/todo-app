package repository

import "context"

func (s *PostgresStorage) UpdateProfile(ctx context.Context, userID int, username, email string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `UPDATE users SET username = $1 WHERE id = $2`, username, userID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE user_passwords SET email = $1 WHERE user_id = $2`, email, userID); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *PostgresStorage) UpdatePasswordHash(ctx context.Context, userID int, passwordHash string) error {
	res, err := s.DB.ExecContext(ctx, `UPDATE user_passwords SET password_hash = $1 WHERE user_id = $2`, passwordHash, userID)
	if err != nil {
		return err
	}
	return checkAffected(res, "пользователь не найден")
}
