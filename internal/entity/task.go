package entity

import "time"

type Task struct {
	ID          int
	Title       string
	Done        bool
	DeletedAt   *time.Time
	UserID      int
	UserTaskId  int
	Description string
}
