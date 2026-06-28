package httpHandler

// import (
// 	"log/slog"
// 	"net/http"
// 	"strconv"
// )

// func (h *TaskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.PathValue("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		h.Logger.Warn("incorrect id", slog.String("id", idStr))
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	err = h.UC.DeleteTask(r.Context(), id)
// 	if err != nil {
// 		h.Logger.Error("server can't delete task with id", slog.Int("id", id))
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// }
