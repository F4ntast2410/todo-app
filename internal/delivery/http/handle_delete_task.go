package httpHandler

import (
	"log/slog"
	"net/http"
	"proj/internal/delivery/http/middleware"
	"strconv"
)

func (h *TaskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	err = h.TaskUC.DeleteTask(r.Context(), taskID)
	if err != nil {
		h.Logger.Error("server can't delete task with id", slog.Int("id", id))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
