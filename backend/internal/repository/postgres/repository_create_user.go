package repository

import (
	"context"
)

func (s *PostgresStorage) CreateUserWeb(ctx context.Context, email string, passwordHash string, username string) error {
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

	err = tx.GetContext(ctx, &lastInsertID, queryUsers, username)
	if err != nil {
		return err
	}

	// 3. Вставляем в таблицу user_passwords (БЕЗ username, его там нет!)
	queryCreds := `INSERT INTO user_passwords (user_id, email, password_hash) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, queryCreds, lastInsertID, email, passwordHash)
	if err != nil {
		return err
	}

	// Записываем полученный ID обратно в структуру, чтобы usecase знал его

	// 4. Применяем изменения, если всё прошло гладко
	return tx.Commit()
}

func (s *PostgresStorage) CreateUserTg(ctx context.Context, ID int64, username string) error {
	// 1. Открываем транзакцию
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 2. Создаем запись в users
	var lastInsertID int
	queryUsers := `INSERT INTO users (username) VALUES ($1) RETURNING id`

	err = tx.GetContext(ctx, &lastInsertID, queryUsers, username)
	if err != nil {
		return err
	}

	// 3. Создаем запись в user_telegram
	queryTg := `INSERT INTO user_telegram (tg_id, user_id, username) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, queryTg, ID, lastInsertID, username)
	if err != nil {
		return err
	}

	return tx.Commit()
}
