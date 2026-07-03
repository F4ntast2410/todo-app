package httpHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	dto "proj/internal/delivery/http/entityRequest"
	"time"
)

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken, err := h.UserUC.LoginUserWeb(r.Context(), req.Email, req.Password)
	if err != nil {
		h.Logger.Warn("invalid credentials or login failed", slog.String("email", req.Email))
		w.WriteHeader(http.StatusUnauthorized) // 401 Unauthorized
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken, // Тот самый UUID сессии из БД
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",                  // Доступно для всех роутов
		HttpOnly: true,                 // Защита от XSS (JS не сможет прочитать куку)
		Secure:   false,                // Поставь true, если будешь использовать HTTPS
		SameSite: http.SameSiteLaxMode, // Защита от CSRF атак
	})

	w.WriteHeader(http.StatusOK)
}
