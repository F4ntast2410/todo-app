package httpHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	dto "proj/internal/delivery/http/entityRequest"
	"proj/internal/delivery/http/middleware"
	"proj/internal/usecase"
)

func (h *UserHandler) GetTelegramLinkHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	link, err := h.UserUC.GetTelegramLink(r.Context(), userID)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		h.Logger.Debug("telegram not linked for user", slog.String("error", err.Error()))
		json.NewEncoder(w).Encode(dto.TelegramLinkResponse{Linked: false})
		return
	}
	json.NewEncoder(w).Encode(dto.TelegramLinkResponse{Linked: true, Username: link.Username})
}

func (h *UserHandler) LinkTelegramHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var req dto.TelegramAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to encode json body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := usecase.VerifyTelegramAuth(req.ToVerifyMap(), h.BotToken); err != nil {
		h.Logger.Warn("telegram auth verification failed", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := h.UserUC.LinkTelegram(r.Context(), userID, req.ID, req.Username); err != nil {
		h.Logger.Warn("failed to link telegram account", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Этот Telegram-аккаунт уже привязан к другому пользователю"})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) TelegramLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.TelegramAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := usecase.VerifyTelegramAuth(req.ToVerifyMap(), h.BotToken); err != nil {
		h.Logger.Warn("telegram auth verification failed", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	loginResult, err := h.UserUC.LoginByTelegram(r.Context(), req.ID, req.Username)
	if err != nil {
		h.Logger.Error("failed telegram login", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    loginResult.SessionToken,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusOK)
}
