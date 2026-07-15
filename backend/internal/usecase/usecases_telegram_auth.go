package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"proj/internal/entity"
	"sort"
	"strings"
)

// VerifyTelegramAuth проверяет hash, который Telegram Login Widget присылает вместе с данными.
// Алгоритм из документации Telegram: https://core.telegram.org/widgets/login#checking-authorization
func VerifyTelegramAuth(data map[string]string, botToken string) error {
	receivedHash := data["hash"]
	if receivedHash == "" {
		return errors.New("hash отсутствует")
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		if k != "hash" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+data[k])
	}
	checkString := strings.Join(pairs, "\n")

	secretKey := sha256.Sum256([]byte(botToken))
	mac := hmac.New(sha256.New, secretKey[:])
	mac.Write([]byte(checkString))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expectedHash), []byte(receivedHash)) {
		return errors.New("подпись Telegram не совпадает")
	}
	return nil
}

func (uc *UserUsecaseImpl) LinkTelegram(ctx context.Context, userID int, tgID int64, tgUsername string) error {
	return uc.UserRepo.LinkTelegramAccount(ctx, userID, tgID, tgUsername)
}

func (uc *UserUsecaseImpl) LoginByTelegram(ctx context.Context, tgID int64, tgUsername string) (*entity.LoginResult, error) {
	userID, err := uc.UserRepo.FindByIdTg(ctx, tgID)
	if err != nil {
		if err := uc.UserRepo.CreateUserTg(ctx, tgID, tgUsername); err != nil {
			return nil, err
		}
		userID, err = uc.UserRepo.FindByIdTg(ctx, tgID)
		if err != nil {
			return nil, err
		}
	}
	sessionID, err := uc.UserRepo.CreateSession(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &entity.LoginResult{SessionToken: sessionID, Requires2FA: false}, nil
}

func (uc *UserUsecaseImpl) GetTelegramLink(ctx context.Context, userID int) (*entity.TelegramLink, error) {
	return uc.UserRepo.GetTelegramLinkByUserID(ctx, userID)
}
