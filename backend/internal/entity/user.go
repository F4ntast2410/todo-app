package entity

import "time"

type User struct {
	ID        int
	Username  string
	CreatedAt time.Time
}

type UserWeb struct {
	UserID       int
	Email        string
	Username     string
	PasswordHash string
}

type UserTg struct {
	ID       int64
	Username string
	UserID   int
}
