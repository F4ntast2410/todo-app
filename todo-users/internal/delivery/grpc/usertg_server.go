package grpc

import (
	"context"
	"errors"
	"proj/internal/entity"
	customErrors "proj/internal/errors"

	"todo-proto/usertg"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserTgUsecase interface {
	RegisterUserTg(ctx context.Context, ID int64, username string) error
	GetUserByTgID(ctx context.Context, userID int64, username string) (*entity.UserTg, error)
}

type UserTgServer struct {
	usertg.UnimplementedUserTgServiceServer
	UserUC UserTgUsecase
}

func NewUserTgServer(userUC UserTgUsecase) *UserTgServer {
	return &UserTgServer{UserUC: userUC}
}

func (s *UserTgServer) RegisterUserTg(ctx context.Context, req *usertg.RegisterUserTgRequest) (*usertg.RegisterUserTgResponse, error) {
	err := s.UserUC.RegisterUserTg(ctx, req.TgId, req.Username)
	if err != nil {
		if errors.Is(err, customErrors.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &usertg.RegisterUserTgResponse{}, nil
}

func (s *UserTgServer) GetUserByTgID(ctx context.Context, req *usertg.GetUserByTgIDRequest) (*usertg.UserTg, error) {
	user, err := s.UserUC.GetUserByTgID(ctx, req.TgId, req.Username)
	if err != nil {
		return nil, err
	}
	return &usertg.UserTg{
		Id:       int64(user.ID),
		Username: user.Username,
		UserId:   int64(user.UserID),
	}, nil
}
