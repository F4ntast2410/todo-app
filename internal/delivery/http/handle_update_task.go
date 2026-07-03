package httpHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	dto "proj/internal/delivery/http/entityRequest"
	"proj/internal/delivery/http/middleware"
	"strconv"
)

func (h *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Warn("incorrect id", slog.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	taskID, err := h.TaskUC.GetTaskByUserID(r.Context(), userID, id)
	if err != nil {
		h.Logger.Warn("task does not exist", slog.Int("user_id", userID), slog.Int("user_task_id", id))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Description != nil {
		err := h.TaskUC.UpdateDescription(r.Context(), taskID, *req.Description)
		if err != nil {
			h.Logger.Error("server can't update task with id", slog.Int("id", taskID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if req.Done != nil {
		err := h.TaskUC.MarkAsDone(r.Context(), taskID, *req.Done)
		if err != nil {
			h.Logger.Error("server can't update task with id", slog.Int("id", taskID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
