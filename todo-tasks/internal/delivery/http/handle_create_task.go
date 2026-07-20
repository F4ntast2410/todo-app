package httpHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	dto "proj/internal/delivery/http/entityRequest"
	"proj/internal/delivery/http/middleware"
)

func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := h.TaskUC.CreateTask(r.Context(), *req.Title, *req.Description, userID)
	if err != nil {
		h.Logger.Error("failed to create new task", slog.String("error", err.Error()), slog.String("attempted_title", *req.Title))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
