package entity

import "time"

type Task struct {
	ID         int        `db:"id" json:"id"`
	Title      string     `db:"title" json:"title"`
	Done       bool       `db:"done" json:"done"`
	DeletedAt  *time.Time `db:"deleted_at" json:"-"`
	UserID     int        `db:"user_id"` // Связующее поле!
	UserTaskId int        `db:"user_task_id"`
}
