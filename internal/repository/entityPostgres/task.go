package pgmodel

import (
	"proj/internal/entity"
	"time"
)

type Task struct {
	ID          int        `db:"id"`
	Title       string     `db:"title"`
	Done        bool       `db:"done"`
	DeletedAt   *time.Time `db:"deleted_at"`
	UserID      int        `db:"user_id"` // Связующее поле!
	UserTaskId  int        `db:"user_task_id"`
	Description string     `db:"description"`
}

func (t Task) ToTask() entity.Task {
	return entity.Task{
		Title:       t.Title,
		UserID:      t.UserID,
		Done:        t.Done,
		Description: t.Description,
	}
}
