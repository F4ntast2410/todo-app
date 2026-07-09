package dto

import (
	"proj/internal/entity"
	"time"
)

type CreateTaskRequest struct {
	Title       *string    `json:"title"`
	Done        *bool      `json:"done"`
	Description *string    `json:"description"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	UserTaskId  int        `json:"user_task_id"`
}
type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (r CreateTaskRequest) ToEntity(userID int) entity.Task {
	return entity.Task{
		Title:       *r.Title,
		UserID:      userID,
		Description: *r.Description,
	}
}
func (r *CreateTaskRequest) ToRequest(t *entity.Task) {
	r.Title = &t.Title
	r.Done = &t.Done
	r.Description = &t.Description
	r.DeletedAt = t.DeletedAt
	r.UserTaskId = t.UserTaskId
}
