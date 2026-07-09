package usecase

import "context"

func (uc *UserUsecaseImpl) GetUserIDBySession(ctx context.Context, sessionID string) (int, error) {
	return uc.UserRepo.FindUserIDBySession(ctx, sessionID)
}
