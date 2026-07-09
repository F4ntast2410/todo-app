package dto

import (
	"proj/internal/entity"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"-"`
}
type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserID   int    `json:"user_id"`
}
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` // сырой пароль, только для входа
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"` // сырой пароль, только для входа
}

type UserTg struct {
	ID       int64  `json:"tg_id"`
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
}

func (u *UserRequest) ToRequest(e *entity.UserWeb) {
	u.Email = e.Email
	u.UserID = e.UserID
	u.Username = e.Username
}
