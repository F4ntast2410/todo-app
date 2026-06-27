package entity

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"-"`
}

type UserWeb struct {
	UserID       int    `json:"user_id"`
	Email        string `son:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type UserTg struct {
	ID       int64  `json:"tg_id"`
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
}
