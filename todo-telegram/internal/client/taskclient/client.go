package taskclient

import (
	"context"
	"proj/internal/entity"
	"time"

	"todo-proto/task"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func toEntityTask(t *task.Task) *entity.Task {
	var deletedAtPointer *time.Time
	if t.DeletedAt != nil {
		deletedAt := t.DeletedAt.AsTime()
		deletedAtPointer = &deletedAt
	}
	return &entity.Task{
		ID:          int(t.Id),
		Title:       t.Title,
		Done:        t.Done,
		UserID:      int(t.UserId),
		UserTaskId:  int(t.UserTaskId),
		Description: t.Description,
		DeletedAt:   deletedAtPointer,
	}
}

type Client struct {
	taskClient task.TaskServiceClient
	conn       *grpc.ClientConn
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		taskClient: task.NewTaskServiceClient(conn),
		conn:       conn,
	}, nil
}

func (c *Client) CreateTask(ctx context.Context, title string, description string, userID int) (*entity.Task, error) {
	t, err := c.taskClient.CreateTask(ctx, &task.CreateTaskRequest{Title: title, Description: description, UserId: int64(userID)})
	if err != nil {
		return nil, err
	}
	return toEntityTask(t), nil
}

func (c *Client) GetTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error) {
	tasks, err := c.taskClient.GetTasksByUserID(ctx, &task.GetTasksByUserIDRequest{UserId: int64(userID)})
	if err != nil {
		return nil, err
	}
	entityTasks := make([]entity.Task, 0, len(tasks.Tasks))
	for _, t := range tasks.Tasks {
		entityTasks = append(entityTasks, *toEntityTask(t))
	}
	return entityTasks, nil
}

func (c *Client) GetRemovedTasksByUserID(ctx context.Context, userID int) ([]entity.Task, error) {
	tasks, err := c.taskClient.GetRemovedTasksByUserID(ctx, &task.GetTasksByUserIDRequest{UserId: int64(userID)})
	if err != nil {
		return nil, err
	}
	entityTasks := make([]entity.Task, 0, len(tasks.Tasks))
	for _, t := range tasks.Tasks {
		entityTasks = append(entityTasks, *toEntityTask(t))
	}
	return entityTasks, nil
}

func (c *Client) GetTask(ctx context.Context, taskID int) (*entity.Task, error) {
	t, err := c.taskClient.GetTask(ctx, &task.GetTaskRequest{TaskId: int64(taskID)})
	if err != nil {
		return nil, err
	}
	return toEntityTask(t), nil
}

func (c *Client) DeleteTask(ctx context.Context, id int) error {
	_, err := c.taskClient.DeleteTask(ctx, &task.TaskIDRequest{Id: int64(id)})
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) DeleteForeverTask(ctx context.Context, id int) error {
	_, err := c.taskClient.DeleteForeverTask(ctx, &task.TaskIDRequest{Id: int64(id)})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) RecoverTask(ctx context.Context, id int) error {
	_, err := c.taskClient.RecoverTask(ctx, &task.TaskIDRequest{Id: int64(id)})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateDescription(ctx context.Context, taskID int, newDesc string) error {
	_, err := c.taskClient.UpdateDescription(ctx, &task.UpdateDescriptionRequest{TaskId: int64(taskID), NewDescription: newDesc})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) MarkAsDone(ctx context.Context, id int, status bool) error {
	_, err := c.taskClient.MarkAsDone(ctx, &task.MarkAsDoneRequest{Id: int64(id), Status: status})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
