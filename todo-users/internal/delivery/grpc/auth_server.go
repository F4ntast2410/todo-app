package grpc

import (
	"context"

	"todo-proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUsecase interface {
	GetUserIDBySession(ctx context.Context, sessionID string) (int, error)
}

// AuthServer реализует интерфейс auth.AuthServiceServer, сгенерированный из .proto
type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	UserUC AuthUsecase
}

func NewAuthServer(userUC AuthUsecase) *AuthServer {
	return &AuthServer{UserUC: userUC}
}

func (s *AuthServer) ValidateSession(ctx context.Context, req *auth.ValidateSessionRequest) (*auth.ValidateSessionResponse, error) {
	userID, err := s.UserUC.GetUserIDBySession(ctx, req.SessionToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired session")
	}
	return &auth.ValidateSessionResponse{UserId: int64(userID)}, nil
}
