package authclient

import (
	"context"

	"todo-proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client реализует тот же интерфейс, что middleware.AuthUsecase ожидает —
// GetUserIDBySession(ctx, sessionID) (int, error). Middleware даже не узнает,
// что теперь это уходит по сети, а не в локальную БД.
type Client struct {
	authClient auth.AuthServiceClient
	conn       *grpc.ClientConn
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		authClient: auth.NewAuthServiceClient(conn),
		conn:       conn,
	}, nil
}

func (c *Client) GetUserIDBySession(ctx context.Context, sessionID string) (int, error) {
	resp, err := c.authClient.ValidateSession(ctx, &auth.ValidateSessionRequest{SessionToken: sessionID})
	if err != nil {
		return 0, err
	}
	return int(resp.UserId), nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
