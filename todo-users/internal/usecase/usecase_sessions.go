package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// generateSessionToken создаёт криптографически случайный токен (уходит клиенту)
// и его SHA-256-хэш (хранится в БД). Раздельно — чтобы сырой токен нигде,
// кроме куки в браузере пользователя, не оседал.
func generateSessionToken() (rawToken string, tokenHash string, err error) {
	raw := make([]byte, 32) // 256 бит энтропии
	if _, err := rand.Read(raw); err != nil {
		return "", "", err
	}
	rawToken = hex.EncodeToString(raw)
	tokenHash = hashToken(rawToken)
	return rawToken, tokenHash, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
func (uc *UserUsecaseImpl) createSessionForUser(ctx context.Context, userID int) (string, error) {
	rawToken, tokenHash, err := generateSessionToken()
	if err != nil {
		return "", err
	}
	if err := uc.AuthRepo.SaveSessionToken(ctx, tokenHash, userID, 24*time.Hour); err != nil {
		return "", err
	}
	err = uc.UserRepo.CreateSession(ctx, userID, tokenHash)
	return rawToken, err
}
func (uc *UserUsecaseImpl) GetUserIDBySession(ctx context.Context, sessionID string) (int, error) {
	return uc.AuthRepo.GetUserID(ctx, hashToken(sessionID))
}
