package httpHandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	dto "proj/internal/delivery/http/entityRequest"
	customErrors "proj/internal/errors"
)

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.UserUC.RegisterUserWeb(r.Context(), req.Email, req.Username, req.Password)
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
}
