package usecase

import (
	"context"
	"proj/internal/entity"
	customErrors "proj/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

func (uc *UserUsecaseImpl) RegisterUserWebSend2FA(ctx context.Context, email entity.VerificationEmail) (*entity.AuthResult, error) {
	exists, err := uc.UserRepo.ExistsWeb(ctx, email)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = uc.Send2FACode(ctx, email)
		if err != nil {
			return nil, err
		}
		return &entity.AuthResult{Requires2FA: true}, nil
	}
	return nil, customErrors.ErrUserAlreadyExists

}
func (uc *UserUsecaseImpl) RegisterUserWeb(ctx context.Context, email entity.VerificationEmail, username, password, inputCode string) (*entity.AuthResult, error) {
	err := uc.Verify2FACode(ctx, email, inputCode)
	if err != nil {
		return nil, err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	userID, err := uc.UserRepo.CreateUserWeb(ctx, email, string(passwordHash), username)
	if err != nil {
		return nil, err
	}
	sessionID, err := uc.createSessionForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &entity.AuthResult{
		SessionToken: sessionID,
		Requires2FA:  false,
	}, nil
}

func (uc *UserUsecaseImpl) RegisterUserTg(ctx context.Context, ID int64, username string) error {
	exists, err := uc.UserRepo.ExistsTg(ctx, ID)
	if err != nil {
		return err
	}
	if !exists {
		err := uc.UserRepo.CreateUserTg(ctx, ID, username)
		if err != nil {
			return err
		}
		return nil
	} else {
		return customErrors.ErrUserAlreadyExists
	}
}
