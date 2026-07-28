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
)

func (uc UserUsecaseImpl) Send2FACode(ctx context.Context, email entity.VerificationEmail) error {
	unique_code, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return err
	}
	unique_code_str := fmt.Sprintf("%d", unique_code.Int64()+100000)
	log.Printf("Code %s", unique_code_str)
	go func() {
		if err := uc.Mailer.Send2FACode(string(email), unique_code_str); err != nil {
			log.Printf("failed to send 2FA email to %s: %v", email, err)
		}
	}()
	if err = uc.AuthRepo.Save2FA(ctx, email, unique_code_str, 5*time.Minute); err != nil {
		return err
	}
	return nil
}
func (uc UserUsecaseImpl) Verify2FACode(ctx context.Context, email entity.VerificationEmail, inputCode string) error {
	savedCode, err := uc.AuthRepo.Get2FA(ctx, email)
	if err != nil {
		return err
	}
	if savedCode != inputCode {
		return customErrors.ErrInvalid2FACode
	}

	err = uc.AuthRepo.Delete2FA(ctx, email)
	if err != nil {
		return err
	}
	return nil
}
