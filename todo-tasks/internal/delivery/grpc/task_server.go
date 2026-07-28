package grpc

import (
	"context"
	"proj/internal/entity"

	"todo-proto/task"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoTask(t *entity.Task) *task.Task {
	var deletedAt *timestamppb.Timestamp
	if t.DeletedAt != nil {
		deletedAt = timestamppb.New(*t.DeletedAt)
	}
	return &task.Task{
		Id:          int64(t.ID),
		Title:       t.Title,
		Done:        t.Done,
		UserId:      int64(t.UserID),
		UserTaskId:  int64(t.UserTaskId),
		Description: t.Description,
		DeletedAt:   deletedAt,
	}
}

type TaskUsecase interface {
	CreateTask(ctx context.Context, title string, description string, userID int) (*entity.Task, error)
	GetTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error)
	GetRemovedTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error)
	GetTask(ctx context.Context, taskID int) (*entity.Task, error)
	DeleteTask(ctx context.Context, id int) error
	DeleteForeverTask(ctx context.Context, id int) error
	RecoverTask(ctx context.Context, id int) error
	UpdateDescription(ctx context.Context, taskID int, newDesc string) error
	MarkAsDone(ctx context.Context, id int, status bool) error
}

type TaskServer struct {
	task.UnimplementedTaskServiceServer
	TaskUC TaskUsecase
}

func NewTaskServer(taskUC TaskUsecase) *TaskServer {
	return &TaskServer{TaskUC: taskUC}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *task.CreateTaskRequest) (*task.Task, error) {
	t, err := s.TaskUC.CreateTask(ctx, req.Title, req.Description, int(req.UserId))
	if err != nil {
		return nil, err
	}
	return toProtoTask(t), nil
}

func (s *TaskServer) GetTasksByUserID(ctx context.Context, req *task.GetTasksByUserIDRequest) (*task.TaskList, error) {
	tasks, err := s.TaskUC.GetTasksByUserID(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	protoTasks := make([]*task.Task, 0, len(tasks))
	for _, t := range tasks {
		protoTasks = append(protoTasks, toProtoTask(&t))
	}
	return &task.TaskList{
		Tasks: protoTasks,
	}, nil
}

func (s *TaskServer) GetRemovedTasksByUserID(ctx context.Context, req *task.GetTasksByUserIDRequest) (*task.TaskList, error) {
	tasks, err := s.TaskUC.GetRemovedTasksByUserID(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	protoTasks := make([]*task.Task, 0, len(tasks))
	for _, t := range tasks {
		protoTasks = append(protoTasks, toProtoTask(&t))
	}
	return &task.TaskList{
		Tasks: protoTasks,
	}, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *task.GetTaskRequest) (*task.Task, error) {
	t, err := s.TaskUC.GetTask(ctx, int(req.TaskId))
	if err != nil {
		return nil, err
	}
	return toProtoTask(t), nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *task.TaskIDRequest) (*task.Empty, error) {
	err := s.TaskUC.DeleteTask(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &task.Empty{}, nil
}
func (s *TaskServer) DeleteForeverTask(ctx context.Context, req *task.TaskIDRequest) (*task.Empty, error) {
	err := s.TaskUC.DeleteForeverTask(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &task.Empty{}, nil
}

func (s *TaskServer) RecoverTask(ctx context.Context, req *task.TaskIDRequest) (*task.Empty, error) {
	err := s.TaskUC.RecoverTask(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &task.Empty{}, nil
}

func (s *TaskServer) UpdateDescription(ctx context.Context, req *task.UpdateDescriptionRequest) (*task.Empty, error) {
	err := s.TaskUC.UpdateDescription(ctx, int(req.TaskId), req.NewDescription)
	if err != nil {
		return nil, err
	}
	return &task.Empty{}, nil
}

func (s *TaskServer) MarkAsDone(ctx context.Context, req *task.MarkAsDoneRequest) (*task.Empty, error) {
	err := s.TaskUC.MarkAsDone(ctx, int(req.Id), req.Status)
	if err != nil {
		return nil, err
	}
	return &task.Empty{}, nil
}
