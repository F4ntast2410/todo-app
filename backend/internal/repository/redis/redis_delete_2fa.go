package cache

import (
	"context"
)

func (c RedisVerificationRepository) Delete2FA(ctx context.Context, email string) error {
	err := c.Client.Del(ctx, "2fa:email:"+email).Err()
	if err != nil {
		return err
	}
	return nil
}
