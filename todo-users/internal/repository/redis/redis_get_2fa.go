package cache

import (
	"context"
	"fmt"
	"proj/internal/entity"
	customErrors "proj/internal/errors"

	"github.com/redis/go-redis/v9"
)

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
