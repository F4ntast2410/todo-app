package httpHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	dto "proj/internal/delivery/http/entityRequest"
	"strconv"
)

func (h *TaskHandler) TrashListHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Warn("incorrect id", slog.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tasks, err := h.UC.GetRemovedTasksByUserID(r.Context(), id)
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req_tasks)
}
