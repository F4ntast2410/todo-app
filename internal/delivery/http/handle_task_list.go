package httpHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	dto "proj/internal/delivery/http/entityRequest"
	"proj/internal/delivery/http/middleware"
)

func (h *TaskHandler) TaskListHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tasks, err := h.TaskUC.GetTasksByUserID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Logger.Error("server can't get tasks", slog.String("error", err.Error()))
		return
	}
	req_tasks := make([]dto.CreateTaskRequest, len(tasks))
	for i, t := range tasks {
		req_tasks[i].ToRequest(&t)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req_tasks)
}
