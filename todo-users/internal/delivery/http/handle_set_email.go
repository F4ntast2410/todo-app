package httpHandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	dto "proj/internal/delivery/http/entityRequest"
	"proj/internal/delivery/http/middleware"
	"proj/internal/entity"
	customErrors "proj/internal/errors"
)

func (h *UserHandler) SetEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestSetEmail
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.UserUC.SetEmail(r.Context(), entity.VerificationEmail(req.Email))
	if err != nil {
		if errors.Is(err, customErrors.ErrUserAlreadyExists) {
			h.Logger.Debug("user already exists", slog.String("error", err.Error()))
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Почта уже занята"})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.Logger.Error("failed to set email", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Logger.Info("success")
	w.WriteHeader(http.StatusOK)
}
func (h *UserHandler) EmailSetVerify2FAHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.Request2FA
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.UserUC.Verify2FACode(r.Context(), entity.VerificationEmail(req.Email), req.InputCode)
	if err != nil {
		if errors.Is(err, customErrors.ErrInvalid2FACode) {
			w.WriteHeader(http.StatusBadRequest) // Выставляем HTTP статус 400
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Error: "Неверный или истекший код подтверждения",
			})
			return
		}
		h.Logger.Error("failed to verify code", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
func (h *UserHandler) SetEmailPasswordHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var req dto.SetEmailPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.UserUC.SetEmailPassword(r.Context(), entity.VerificationEmail(req.Email), req.Password, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
