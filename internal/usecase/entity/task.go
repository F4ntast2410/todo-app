package entity

import "time"

type Task struct {
	ID         int        `json:"id"`
	Title      string     `json:"title"`
	Done       bool       `json:"done"`
	DeletedAt  *time.Time `json:"-"`
	UserID     int        `json:"-"` // Связующее поле!
	UserTaskId int        `json:"-"`
}
