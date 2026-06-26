package repository

import (
	"proj/internal/entity"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Task = entity.Task

type PostgresStorage struct {
	DB *sqlx.DB
}
