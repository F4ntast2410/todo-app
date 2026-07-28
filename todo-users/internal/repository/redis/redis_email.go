package cache

import (
	"context"
	"fmt"
	"proj/internal/entity"
	customErrors "proj/internal/errors"
	"time"

	"github.com/redis/go-redis/v9"
)

func (c RedisVerificationRepository) Delete2FA(ctx context.Context, email entity.VerificationEmail) error {
	err := c.Client.Del(ctx, fmt.Sprintf("2fa:email:%s", email)).Err()
	if err != nil {
		return err
	}
	return nil
}
func (c RedisVerificationRepository) Get2FA(ctx context.Context, email entity.VerificationEmail) (string, error) {
	result, err := c.Client.Get(ctx, fmt.Sprintf("2fa:email:%s", email)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", customErrors.ErrCacheValueNotExists
		}
		return "", err
	}
	return result, err
}
func (c RedisVerificationRepository) Save2FA(ctx context.Context, email entity.VerificationEmail, unique_code string, time_duration time.Duration) error {
	err := c.Client.Set(ctx, fmt.Sprintf("2fa:email:%s", email), unique_code, time_duration).Err()
	if err != nil {
		return err
	}
	return nil
}
