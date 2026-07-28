package cache

import (
	"context"
	"fmt"
	customErrors "proj/internal/errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func (c RedisVerificationRepository) DeleteSessionToken(ctx context.Context, sessionID string) error {
	err := c.Client.Del(ctx, fmt.Sprintf("session:session_token:%s", sessionID)).Err()
	if err != nil {
		return err
	}
	return nil
}
func (c RedisVerificationRepository) GetUserID(ctx context.Context, sessionID string) (int, error) {
	result, err := c.Client.Get(ctx, fmt.Sprintf("session:session_token:%s", sessionID)).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, customErrors.ErrCacheValueNotExists
		}
		return 0, err
	}
	userID, err := strconv.Atoi(result)
	return userID, err
}
func (c RedisVerificationRepository) SaveSessionToken(ctx context.Context, sessionID string, userID int, time_duration time.Duration) error {
	err := c.Client.Set(ctx, fmt.Sprintf("session:session_token:%s", sessionID), userID, time_duration).Err()
	if err != nil {
		return err
	}
	return nil
}
