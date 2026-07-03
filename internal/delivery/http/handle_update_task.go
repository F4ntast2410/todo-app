package httpHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	dto "proj/internal/delivery/http/entityRequest"
	"strconv"
)

func (h *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Warn("incorrect id", slog.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Description != nil {
		err := h.UC.UpdateDescription(r.Context(), id, *req.Description)
		if err != nil {
			h.Logger.Error("server can't update task with id", slog.Int("id", id))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if req.Done != nil {
		err := h.UC.MarkAsDone(r.Context(), req.ID, *req.Done)
		if err != nil {
			h.Logger.Error("server can't update task with id", slog.Int("id", req.ID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
