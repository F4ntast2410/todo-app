package usecase

import (
	"context"
	"time"
)

func (uc UserUsecaseImpl) Save2FA(ctx context.Context, email string, unique_code string, time_duration time.Duration) error {
	return uc.AuthRepo.Save2FA(ctx, email, unique_code, time_duration)
}

func (uc UserUsecaseImpl) Get2FA(ctx context.Context, email string) (string, error) {
	return uc.AuthRepo.Get2FA(ctx, email)
}

func (uc UserUsecaseImpl) Delete2FA(ctx context.Context, email string) error {
	return uc.AuthRepo.Delete2FA(ctx, email)
}
