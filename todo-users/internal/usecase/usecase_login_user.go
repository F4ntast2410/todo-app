package usecase

import (
	"context"
	"proj/internal/entity"

	"golang.org/x/crypto/bcrypt"
)

func (uc *UserUsecaseImpl) LoginUserWeb(ctx context.Context, email entity.VerificationEmail, password string) (*entity.AuthResult, error) {
	user, err := uc.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}
	err = uc.Send2FACode(ctx, email)
	return &entity.AuthResult{Requires2FA: true}, err
}
func (uc *UserUsecaseImpl) Verify2FA(ctx context.Context, email entity.VerificationEmail, inputCode string) (*entity.AuthResult, error) {
	user, err := uc.UserRepo.FindByEmail(ctx, entity.VerificationEmail(email))
	if err != nil {
		return nil, err
	}
	err = uc.Verify2FACode(ctx, email, inputCode)
	if err != nil {
		return nil, err
	}
	sessionID, err := uc.createSessionForUser(ctx, user.UserID)
	if err != nil {
		return nil, err
	}
	return &entity.AuthResult{
		SessionToken: sessionID,
		Requires2FA:  false,
	}, nil
}
