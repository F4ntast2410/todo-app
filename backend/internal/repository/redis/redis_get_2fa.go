package cache

import (
	"context"
	customErrors "proj/internal/errors"

	"github.com/redis/go-redis/v9"
)

func (c RedisVerificationRepository) Get2FA(ctx context.Context, email string) (string, error) {
	result, err := c.Client.Get(ctx, "2fa:email:"+email).Result()
	if err != nil {
		if err == redis.Nil {
			return "", customErrors.ErrCacheValueNotExists
		}
		return "", err
	}
	return result, err
}
