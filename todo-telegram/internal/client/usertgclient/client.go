package usertgclient

import (
	"context"
	"proj/internal/entity"
	customErrors "proj/internal/errors"

	"todo-proto/usertg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	userTgClient usertg.UserTgServiceClient
	conn         *grpc.ClientConn
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		userTgClient: usertg.NewUserTgServiceClient(conn),
		conn:         conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) RegisterUserTg(ctx context.Context, ID int64, username string) error {
	_, err := c.userTgClient.RegisterUserTg(ctx, &usertg.RegisterUserTgRequest{TgId: ID, Username: username})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return customErrors.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (c *Client) GetUserByTgID(ctx context.Context, userID int64, username string) (*entity.UserTg, error) {
	user, err := c.userTgClient.GetUserByTgID(ctx, &usertg.GetUserByTgIDRequest{TgId: userID, Username: username})
	if err != nil {
		return nil, err
	}
	return &entity.UserTg{ID: user.Id, Username: user.Username, UserID: int(user.UserId)}, nil
}
