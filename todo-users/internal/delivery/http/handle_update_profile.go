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

func (h *UserHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var req dto.UpdateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Email != nil {
		err = h.UserUC.UpdateEmail(r.Context(), entity.VerificationEmail(*req.NewEmail), entity.VerificationEmail(*req.Email))
		if err != nil {
			if errors.Is(err, customErrors.ErrUserAlreadyExists) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}
	if req.Username != nil {
		err = h.UserUC.UpdateUsername(r.Context(), userID, *req.Username)
	}
	if err != nil {
		h.Logger.Error("failed to update profile", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (h *UserHandler) UpdateEmailVerifyHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var req dto.Request2FA
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.UserUC.UpdateEmailVerify(r.Context(), userID, entity.VerificationEmail(req.NewEmail), entity.VerificationEmail(req.Email), req.InputCode)
	if err != nil {
		if errors.Is(err, customErrors.ErrInvalid2FACode) {
			w.WriteHeader(http.StatusBadRequest)
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

func (h *UserHandler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.UserUC.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		h.Logger.Warn("failed to change password", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Неверный текущий пароль"})
		return
	}
	w.WriteHeader(http.StatusOK)
}
