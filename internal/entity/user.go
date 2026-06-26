package entity

import "time"

type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at"`
}

type UserWeb struct {
	UserID       int    `db:"user_id" json:"user_id"`
	Email        string `db:"email" json:"email"`
	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"password_hash" json:"-"`
}

type UserTg struct {
	ID       int64  `db:"tg_id"`
	Username string `db:"username"`
	UserID   int    `db:"user_id"`
}
