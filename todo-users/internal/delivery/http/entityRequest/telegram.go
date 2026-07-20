package dto

import (
	"fmt"
	"strconv"
)

// TelegramAuthRequest — то, что присылает Telegram Login Widget (и в поле привязки, и на логин-странице)
type TelegramAuthRequest struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

// ToVerifyMap собирает данные в формате, нужном для проверки HMAC-подписи Telegram.
// Важно: пустые необязательные поля не включаются — Telegram их тоже не отправляет, если они не заданы у пользователя.
func (r TelegramAuthRequest) ToVerifyMap() map[string]string {
	m := map[string]string{
		"id":        strconv.FormatInt(r.ID, 10),
		"auth_date": fmt.Sprintf("%d", r.AuthDate),
		"hash":      r.Hash,
	}
	if r.FirstName != "" {
		m["first_name"] = r.FirstName
	}
	if r.LastName != "" {
		m["last_name"] = r.LastName
	}
	if r.Username != "" {
		m["username"] = r.Username
	}
	if r.PhotoURL != "" {
		m["photo_url"] = r.PhotoURL
	}
	return m
}
