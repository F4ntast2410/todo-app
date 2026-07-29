package repository

import (
	"context"
	"database/sql"
	"proj/internal/entity"
)

func (s *PostgresStorage) FindByIdTg(ctx context.Context, tgID int64) (int, error) {
	var userID int
	query := `SELECT user_id FROM user_telegram WHERE tg_id = $1`
	err := s.DB.GetContext(ctx, &userID, query, tgID)
	return userID, err
}

func (s *PostgresStorage) FindByIdWeb(ctx context.Context, userID int) (*entity.UserWeb, error) {
	var row struct {
		Username     string         `db:"username"`
		Email        sql.NullString `db:"email"`
		PasswordHash sql.NullString `db:"password_hash"`
	}

	query := `SELECT u.username, up.email, up.password_hash 
FROM users u 
LEFT JOIN user_passwords up ON u.id = up.user_id 
WHERE u.id = $1`

	err := s.DB.GetContext(ctx, &row, query, userID)
	if err != nil {
		return nil, err
	}

	user := entity.UserWeb{
		UserID:       userID,
		Username:     row.Username,
		Email:        row.Email.String,
		PasswordHash: row.PasswordHash.String,
		HasPassword:  row.Email.Valid,
	}
	return &user, nil
}
