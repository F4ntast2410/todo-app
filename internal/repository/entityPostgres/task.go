package pgmodel

import "time"

type Task struct {
	ID         int        `db:"id"`
	Title      string     `db:"title"`
	Done       bool       `db:"done"`
	DeletedAt  *time.Time `db:"deleted_at"`
	UserID     int        `db:"user_id"` // Связующее поле!
	UserTaskId int        `db:"user_task_id"`
}
