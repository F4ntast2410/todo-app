package usecase

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

func (uc *UserUsecaseImpl) LoginUserWeb(ctx context.Context, email, password string) (string, error) {
	user, err := uc.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}
	session_id, err := uc.UserRepo.CreateSession(ctx, user.UserID)
	if err != nil {
		return "", err
	}
	return session_id, nil
}
