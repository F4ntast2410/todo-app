package dto

import "proj/internal/entity"

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Done        bool   `json:"done"`
	Description string `json:"description"`
}

func (r CreateTaskRequest) ToEntity(userID int) entity.Task {
	return entity.Task{
		Title:       r.Title,
		UserID:      userID,
		Description: r.Description,
	}
}
func (r *CreateTaskRequest) ToRequest(t *entity.Task) {
	r.Title = t.Title
	r.Done = t.Done
	r.Description = t.Description
}
