package cache

import (
	"context"
	"time"
)

func (c RedisVerificationRepository) Save2FA(ctx context.Context, email string, unique_code string, time_duration time.Duration) error {
	err := c.Client.Set(ctx, "2fa:email:"+email, unique_code, time_duration).Err()
	if err != nil {
		return err
	}
	return nil
}
