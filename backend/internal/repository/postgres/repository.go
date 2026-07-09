package repository

import (
	pgmodel "proj/internal/repository/postgres/entityPostgres"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Task = pgmodel.Task

type PostgresStorage struct {
	DB *sqlx.DB
}

func NewPostgressStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{DB: db}
}
