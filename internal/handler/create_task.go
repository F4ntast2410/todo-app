package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
		ID    int    `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t, err := h.UC.CreateTask(r.Context(), req.Title, req.ID)
	if err != nil {
		h.Logger.Error("failed to create new task", slog.String("error", err.Error()), slog.String("attempted_title", req.Title))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}
