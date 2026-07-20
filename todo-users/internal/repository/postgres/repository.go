package repository

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresStorage struct {
	DB *sqlx.DB
}

func NewPostgressStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{DB: db}
}
