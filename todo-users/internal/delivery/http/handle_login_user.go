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

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	loginResult, err := h.UserUC.LoginUserWeb(r.Context(), entity.VerificationEmail(req.Email), req.Password)
	if err != nil {
		h.Logger.Warn("invalid credentials or login failed", slog.String("email", req.Email))
		w.WriteHeader(http.StatusUnauthorized) // 401 Unauthorized
		return
	}
	resp := dto.Responce2FA{Status: loginResult.Requires2FA}
	if resp.Status == true {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) Verify2FA(w http.ResponseWriter, r *http.Request) {
	var req dto.Request2FA
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	loginResult, err := h.UserUC.Verify2FA(r.Context(), entity.VerificationEmail(req.Email), req.InputCode)
	if err != nil {
		if errors.Is(err, customErrors.ErrInvalid2FACode) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Error: "Неверный или истекший код подтверждения",
			})
			return
		}
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
func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Успешный выход"}`))
}
