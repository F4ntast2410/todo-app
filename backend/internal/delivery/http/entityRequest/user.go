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
type UserWeb struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (u *UserRequest) ToRequest(e *entity.UserWeb) {
	u.Email = e.Email
	u.UserID = e.UserID
	u.Username = e.Username
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
type TelegramLinkResponse struct {
	Linked   bool   `json:"linked"`
	Username string `json:"username,omitempty"`
}
