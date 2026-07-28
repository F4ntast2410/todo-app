package httpHandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	dto "proj/internal/delivery/http/entityRequest"
	"proj/internal/entity"
	customErrors "proj/internal/errors"
)

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	registerResult, err := h.UserUC.RegisterUserWebSend2FA(r.Context(), entity.VerificationEmail(req.Email))
	if err != nil {
		if errors.Is(err, customErrors.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict) // 409 — ресурс уже существует
			return
		}
		h.Logger.Error("failed register user", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.Responce2FA{Status: registerResult.Requires2FA})
}
func (h *UserHandler) RegisterVerify2FAHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.Request2FARegister
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	registerResult, err := h.UserUC.RegisterUserWeb(r.Context(), entity.VerificationEmail(req.Email), req.Username, req.Password, req.InputCode)
	if err != nil {
		if errors.Is(err, customErrors.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict) // 409 — ресурс уже существует
			return
		}
		h.Logger.Error("failed register user", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    registerResult.SessionToken,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusCreated)
}
