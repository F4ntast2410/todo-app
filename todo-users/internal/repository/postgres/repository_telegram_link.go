package repository

import (
	"context"
	"proj/internal/entity"
	pgmodel "proj/internal/repository/postgres/entityPostgres"
)

func (s *PostgresStorage) LinkTelegramAccount(ctx context.Context, userID int, tgID int64, username string) error {
	query := `INSERT INTO user_telegram (tg_id, user_id, username) VALUES ($1, $2, $3)`
	_, err := s.DB.ExecContext(ctx, query, tgID, userID, username)
	return err
}

func (s *PostgresStorage) GetTelegramLinkByUserID(ctx context.Context, userID int) (*entity.TelegramLink, error) {
	var link pgmodel.TelegramLink
	query := `SELECT tg_id, username FROM user_telegram WHERE user_id = $1`
	err := s.DB.GetContext(ctx, &link, query, userID)
	if err != nil {
		return nil, err
	}
	link_entity := link.ToEntitiy()
	return &link_entity, nil
}

func (s *PostgresStorage) LinkWebAccount(ctx context.Context, userID int, email entity.VerificationEmail, passwordHash string) error {
	queryCreds := `INSERT INTO user_passwords (user_id, email, password_hash) VALUES ($1, $2, $3)`
	_, err := s.DB.ExecContext(ctx, queryCreds, userID, email, passwordHash)
	return err
}
