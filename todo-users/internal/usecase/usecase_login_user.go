package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"proj/internal/entity"
	customErrors "proj/internal/errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (uc *UserUsecaseImpl) LoginUserWeb(ctx context.Context, email entity.VerificationEmail, password string) (*entity.LoginResult, error) {
	user, err := uc.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}
	unique_code, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return nil, err
	}
	unique_code_str := fmt.Sprintf("%d", unique_code.Int64()+100000)
	go func() {
		if err := uc.Mailer.Send2FACode(string(email), unique_code_str); err != nil {
			log.Printf("failed to send 2FA email to %s: %v", email, err)
		}
	}()
	if err = uc.AuthRepo.Save2FA(ctx, email, unique_code_str, 5*time.Minute); err != nil {
		return nil, err
	}
	return &entity.LoginResult{Requires2FA: true}, nil
}
func (uc *UserUsecaseImpl) Verify2FA(ctx context.Context, email entity.VerificationEmail, inputCode string) (*entity.LoginResult, error) {
	user, err := uc.UserRepo.FindByEmail(ctx, entity.VerificationEmail(email))
	if err != nil {
		return nil, err
	}
	savedCode, err := uc.AuthRepo.Get2FA(ctx, email)
	if err != nil {
		return nil, err
	}
	if savedCode != inputCode {
		return nil, customErrors.ErrInvalid2FACode
	}

	err = uc.AuthRepo.Delete2FA(ctx, email)
	if err != nil {
		return nil, err
	}
	sessionID, err := uc.UserRepo.CreateSession(ctx, user.UserID)
	if err != nil {
		return nil, err
	}
	return &entity.LoginResult{
		SessionToken: sessionID,
		Requires2FA:  false,
	}, nil
}
