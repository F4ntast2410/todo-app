package dto

import "proj/internal/entity"

type CreateTaskRequest struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (r CreateTaskRequest) ToEntity(userID int) entity.Task {
	return entity.Task{
		Title:  r.Title,
		UserID: userID,
	}
}
