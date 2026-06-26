package handler

// import (
// 	"net/http"
// )

// func (h *TaskHandler) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {

// 	tasks, err := h.UC.DefaultGetAllTasks(r.Context())
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		h.Logger.Error("server can't get tasks", slog.String("error", err.Error()))
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(tasks)
// }
