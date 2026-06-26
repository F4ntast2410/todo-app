package repository

import (
	"context"
	"proj/internal/entity"
)

func (s *PostgresStorage) CreateUserWeb(ctx context.Context, user *entity.UserWeb) error {
	// 1. Открываем транзакцию
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	// Если что-то пойдет не так, defer автоматически откатит изменения
	defer tx.Rollback()

	// 2. Вставляем в родительскую таблицу users и получаем сгенерированный ID
	var lastInsertID int
	queryUsers := `INSERT INTO users (username) VALUES ($1) RETURNING id`

	err = tx.GetContext(ctx, &lastInsertID, queryUsers, user.Username)
	if err != nil {
		return err
	}

	// 3. Вставляем в таблицу user_passwords (БЕЗ username, его там нет!)
	queryCreds := `INSERT INTO user_passwords (user_id, email, password_hash) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, queryCreds, lastInsertID, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	// Записываем полученный ID обратно в структуру, чтобы usecase знал его
	user.UserID = lastInsertID

	// 4. Применяем изменения, если всё прошло гладко
	return tx.Commit()
}

func (s *PostgresStorage) CreateUserTg(ctx context.Context, user *entity.UserTg) error {
	// 1. Открываем транзакцию
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 2. Создаем запись в users
	var lastInsertID int
	queryUsers := `INSERT INTO users (username) VALUES ($1) RETURNING id`

	err = tx.GetContext(ctx, &lastInsertID, queryUsers, user.Username)
	if err != nil {
		return err
	}

	// 3. Создаем запись в user_telegram
	queryTg := `INSERT INTO user_telegram (tg_id, user_id, username) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, queryTg, user.ID, lastInsertID, user.Username)
	if err != nil {
		return err
	}

	user.UserID = lastInsertID

	return tx.Commit()
}
